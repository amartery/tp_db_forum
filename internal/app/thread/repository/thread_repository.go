package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/amartery/tp_db_forum/internal/app/forum"
	postModel "github.com/amartery/tp_db_forum/internal/app/post/models"
	"github.com/amartery/tp_db_forum/internal/app/thread/models"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ThreadRepository struct {
	Con *pgxpool.Pool
}

func NewThreadRepository(con *pgxpool.Pool) *ThreadRepository {
	return &ThreadRepository{
		Con: con,
	}
}

func (t *ThreadRepository) CheckForum(slug string) (string, error) {
	err := t.Con.QueryRow(
		context.Background(),
		"SELECT slug FROM forums WHERE slug = $1",
		slug,
	).Scan(&slug)

	if err != nil {
		return "", fmt.Errorf("couldn't get forum with slug '%v'. Error: %w", slug, err)
	}

	return slug, nil
}

func (t *ThreadRepository) CreateThread(thread *models.Thread) error {
	var err error
	thread.Forum, err = t.CheckForum(thread.Forum)

	if err != nil {
		return forum.ErrForumDoesntExists
	}

	query := `INSERT INTO threads (author, created, forum, msg, title, slug)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = t.Con.QueryRow(
		context.Background(),
		query,
		thread.Author, thread.Created, thread.Forum, thread.Message, thread.Title, thread.Slug,
	).Scan(&thread.ID)

	if err != nil {
		return fmt.Errorf("couldn't create thread. Error: %w", err)
	}

	return nil
}

func (t *ThreadRepository) GetThreadsByForumSlug(slug, limit, since, desc string) (*[]models.Thread, error) {
	query := "SELECT author, created, forum, id, msg, slug, title, votes FROM threads WHERE forum = $1"

	args := make([]interface{}, 0, 4)
	args = append(args, slug)

	var operator string
	if desc == "" || desc == "false" {
		operator = ">"
	} else {
		operator = "<"
	}

	if since != "" {
		query += fmt.Sprintf(" AND created %v= $2", operator)
		args = append(args, since)
	}

	if desc == "" || desc == "false" {
		desc = "ASC"
	} else {
		desc = "DESC"
	}
	query += fmt.Sprintf(" ORDER BY created %v", desc)

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	query += fmt.Sprintf(" LIMIT %v", limitInt)

	rows, err := t.Con.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threads := make([]models.Thread, 0, limitInt)
	var thread models.Thread
	for rows.Next() {
		err = rows.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
		if err != nil {
			return nil, err
		}

		threads = append(threads, thread)
	}

	return &threads, nil
}

func (tr *ThreadRepository) CreatePosts(thread models.Thread, posts []postModel.Post) error {
	if posts[0].Parent != 0 {
		var parentThread int
		err := tr.Con.QueryRow(context.Background(),
			"SELECT thread FROM posts WHERE id = $1",
			posts[0].Parent,
		).Scan(&parentThread)

		if err != nil {
			return fmt.Errorf("couldn't get thread id from posts: %w", err)
		}

		if parentThread != thread.ID {
			return forum.ErrWrongParent
		}
	}

	query := `INSERT INTO posts(author, created, forum, msg, parent, thread) VALUES `
	var args []interface{}
	created := strfmt.DateTime(time.Now())

	for i, post := range posts {
		posts[i].Forum = thread.Forum
		posts[i].Thread = thread.ID
		posts[i].Created = created

		query += fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d),",
			i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6,
		)

		args = append(args, post.Author, created, thread.Forum, post.Message, post.Parent, thread.ID)
	}

	query = query[:len(query)-1]
	query += ` RETURNING id`

	rows, err := tr.Con.Query(context.Background(), query, args...)
	if err != nil {
		return fmt.Errorf("couldn't insert posts: %w", err)
	}
	defer rows.Close()

	var idx int
	for rows.Next() {
		err := rows.Scan(&posts[idx].ID)
		if err != nil {
			return fmt.Errorf("couldn't scan post id: %w", err)
		}

		idx++
	}

	return nil

}

func (tr *ThreadRepository) GetPosts(slugOrID string, limit int, order string, since string) ([]postModel.Post, error) {
	var sinceCond string
	if since != "" {
		if order == "DESC" {
			sinceCond = fmt.Sprintf("AND id < %v", since)
		} else {
			sinceCond = fmt.Sprintf("AND id > %v", since)
		}
	}

	threadID, err := strconv.Atoi(slugOrID)
	if err != nil {
		threadID, err = tr.CheckThreadBySlug(slugOrID)
		if err != nil {
			return nil, err
		}
	}

	query := fmt.Sprintf(`SELECT author, created, forum, id, msg, parent, thread FROM posts
	WHERE thread = $1 %v
	ORDER BY id %v`, sinceCond, order)

	if limit != 0 {
		query += fmt.Sprintf(" LIMIT %v", limit)
	}

	rows, err := tr.Con.Query(context.Background(), query, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]postModel.Post, 0, limit)
	post := postModel.Post{}
	for rows.Next() {
		err = rows.Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.Message, &post.Parent, &post.Thread)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (tr *ThreadRepository) GetThreadBySlug(slug string) (*models.Thread, error) {
	query := `SELECT author, created, forum, id, msg, slug, title, votes FROM threads WHERE slug = $1`
	thread := &models.Thread{}
	err := tr.Con.QueryRow(context.Background(), query, slug).Scan(
		&thread.Author,
		&thread.Created,
		&thread.Forum,
		&thread.ID,
		&thread.Message,
		&thread.Slug,
		&thread.Title,
		&thread.Votes)

	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (tr *ThreadRepository) GetThreadByID(id int) (*models.Thread, error) {
	query := `SELECT author, created, forum, id, msg, slug, title, votes FROM threads WHERE id = $1`
	thread := &models.Thread{}
	err := tr.Con.QueryRow(context.Background(), query, id).Scan(
		&thread.Author,
		&thread.Created,
		&thread.Forum,
		&thread.ID,
		&thread.Message,
		&thread.Slug,
		&thread.Title,
		&thread.Votes)

	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (tr *ThreadRepository) GetThreadIDAndForum(slugOrID string) (*models.Thread, error) {
	threadID, err := strconv.Atoi(slugOrID)
	thread := &models.Thread{ID: threadID}
	if err != nil {
		query := `SELECT forum, id FROM threads WHERE slug = $1`
		err = tr.Con.QueryRow(context.Background(), query, slugOrID).Scan(&thread.Forum, &thread.ID)
	} else {
		query := `SELECT forum FROM threads WHERE id = $1`
		err = tr.Con.QueryRow(context.Background(), query, thread.ID).Scan(&thread.Forum)
	}

	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (tr *ThreadRepository) CheckThreadBySlug(slug string) (int, error) {
	var id int
	err := tr.Con.QueryRow(context.Background(),
		"SELECT id FROM threads WHERE slug = $1",
		slug,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (tr *ThreadRepository) CheckThreadByID(id int) (int, error) {
	err := tr.Con.QueryRow(context.Background(),
		"SELECT id FROM threads WHERE id = $1",
		id,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (tr *ThreadRepository) UpdateThread(thread *models.Thread) (*models.Thread, error) {
	if thread.Title == "" || thread.Message == "" {
		oldThread := &models.Thread{}
		var err error

		if thread.ID != 0 {
			oldThread, err = tr.GetThreadByID(thread.ID)
		} else {
			oldThread, err = tr.GetThreadBySlug(*thread.Slug)
		}
		if err != nil {
			return nil, err
		}

		if thread.Title == "" {
			thread.Title = oldThread.Title
		}

		if thread.Message == "" {
			thread.Message = oldThread.Message
		}
	}

	err := tr.Con.QueryRow(context.Background(),
		`UPDATE threads SET title = $1, msg = $2
		WHERE slug = $3 OR id = $4
		RETURNING author, created, forum, id, msg, slug, title`,
		thread.Title, thread.Message, thread.Slug, thread.ID,
	).Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &thread.Slug, &thread.Title)

	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (tr *ThreadRepository) GetPostsTree(slugOrID string, limit int, order string, since string) ([]postModel.Post, error) {
	var desc bool
	if order == "DESC" {
		desc = true
	} else {
		desc = false
	}

	threadID, err := strconv.Atoi(slugOrID)
	if err != nil {
		threadID, err = tr.CheckThreadBySlug(slugOrID)
		if err != nil {
			return nil, err
		}
	}

	///////

	var query string

	if since == "" {
		if desc {
			query = fmt.Sprintf(`SELECT author, created, forum, id, msg, parent, thread FROM posts
				WHERE thread = %d ORDER BY path DESC, id  DESC LIMIT %d;`, threadID, limit)
		} else {
			query = fmt.Sprintf(`SELECT author, created, forum, id, msg, parent, thread FROM posts
				WHERE thread = %d ORDER BY path ASC, id  ASC LIMIT %d;`, threadID, limit)
		}
	} else {
		if desc {
			query = fmt.Sprintf(`SELECT author, created, forum, id, msg, parent, thread FROM posts
				WHERE thread = %d AND path < (SELECT path FROM posts WHERE id = %s)
				ORDER BY path DESC, id  DESC LIMIT %d;`, threadID, since, limit)
		} else {
			query = fmt.Sprintf(`SELECT author, created, forum, id, msg, parent, thread FROM posts
				WHERE thread = %d AND path > (SELECT path FROM posts WHERE id = %s)
				ORDER BY path ASC, id  ASC LIMIT %d;`, threadID, since, limit)
		}
	}

	rows, err := tr.Con.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	posts := make([]postModel.Post, 0)
	var post postModel.Post
	for rows.Next() {
		err = rows.Scan(
			&post.Author, &post.Created, &post.Forum, &post.ID,
			&post.Message, &post.Parent, &post.Thread,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (tr *ThreadRepository) GetPostsParentTree(slugOrID string, limit int, order string, since string) ([]postModel.Post, error) {
	var desc bool
	if order == "DESC" {
		desc = true
	} else {
		desc = false
	}

	threadID, err := strconv.Atoi(slugOrID)
	if err != nil {
		threadID, err = tr.CheckThreadBySlug(slugOrID)
		if err != nil {
			return nil, err
		}
	}

	var query string
	if since == "" {
		if desc {
			query = fmt.Sprintf(`SELECT author, created, forum, id, msg, parent, thread FROM posts
				WHERE path[1] IN (SELECT id FROM posts WHERE thread = %d AND parent = 0 ORDER BY id DESC LIMIT %d)
				ORDER BY path[1] DESC, path, id;`, threadID, limit)
		} else {
			query = fmt.Sprintf(`SELECT author, created, forum, id, msg, parent, thread FROM posts
				WHERE path[1] IN (SELECT id FROM posts WHERE thread = %d AND parent = 0 ORDER BY id LIMIT %d)
				ORDER BY path, id;`, threadID, limit)
		}
	} else {
		if desc {
			query = fmt.Sprintf(`SELECT author, created, forum, id, msg, parent, thread FROM posts
				WHERE path[1] IN (SELECT id FROM posts WHERE thread = %d AND parent = 0 AND path[1] <
				(SELECT path[1] FROM posts WHERE id = %s) ORDER BY id DESC LIMIT %d) ORDER BY path[1] DESC, path, id;`,
				threadID, since, limit)
		} else {
			query = fmt.Sprintf(`SELECT author, created, forum, id, msg, parent, thread FROM posts
				WHERE path[1] IN (SELECT id FROM posts WHERE thread = %d AND parent = 0 AND path[1] >
				(SELECT path[1] FROM posts WHERE id = %s) ORDER BY id ASC LIMIT %d) ORDER BY path, id;`,
				threadID, since, limit)
		}
	}

	rows, err := tr.Con.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := make([]postModel.Post, 0)
	var post postModel.Post
	for rows.Next() {
		err = rows.Scan(
			&post.Author, &post.Created, &post.Forum, &post.ID,
			&post.Message, &post.Parent, &post.Thread,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (tr *ThreadRepository) Vote(vote *models.Vote) (*models.Thread, error) {

	thread := &models.Thread{}
	var err error

	if vote.ID != 0 {
		thread, err = tr.GetThreadByID(vote.ID)
	} else {
		thread, err = tr.GetThreadBySlug(vote.Slug)
	}

	if err != nil {
		return nil, err
	}

	var voteValue int
	err = tr.Con.QueryRow(context.Background(),
		`SELECT vote FROM thread_vote
		WHERE nickname = $1 AND thread_id = $2`,
		vote.Nickname, thread.ID,
	).Scan(&voteValue)

	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	if err == pgx.ErrNoRows {
		_, err = tr.Con.Exec(context.Background(),
			"INSERT INTO thread_vote (nickname, thread_id, vote) VALUES($1, $2, $3)",
			vote.Nickname, thread.ID, vote.Voice,
		)

		if err != nil {
			return nil, err
		}

		thread.Votes += vote.Voice
		return thread, nil
	}

	if voteValue == vote.Voice {
		return thread, nil
	}

	thread.Votes = thread.Votes - voteValue + vote.Voice

	_, err = tr.Con.Exec(context.Background(),
		`UPDATE thread_vote SET vote = $1
		WHERE nickname = $2 AND thread_id = $3`,
		vote.Voice, vote.Nickname, thread.ID,
	)

	if err != nil {
		return nil, err
	}

	return thread, nil
}
