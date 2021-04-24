package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	postModels "github.com/amartery/tp_db_forum/internal/app/post/models"
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

func (t *ThreadRepository) FindThreadBySlug(slug string) (*models.Thread, error) {
	query := `SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE slug = $1`
	thread := &models.Thread{}
	createdTime := &time.Time{}

	err := t.Con.QueryRow(
		context.Background(),
		query,
		slug).Scan(
		&thread.ID,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		createdTime)
	thread.Created = strfmt.DateTime(createdTime.UTC()).String()
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (t *ThreadRepository) CreateThread(thread *models.Thread) (*models.Thread, error) {
	var err error
	if thread.Created != "" {
		query := `INSERT INTO Threads (title, author, forum, message, slug, created)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, title, author, forum, message, votes, slug`

		err = t.Con.QueryRow(
			context.Background(),
			query,
			thread.Title,
			thread.Author,
			thread.Forum,
			thread.Message,
			thread.Slug,
			thread.Created).Scan(
			&thread.ID,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug)
	} else {
		query := `INSERT INTO Threads (title, author, forum, message, slug)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, title, author, forum, message, votes, slug`

		err = t.Con.QueryRow(
			context.Background(),
			query,
			thread.Title,
			thread.Author,
			thread.Forum,
			thread.Message,
			thread.Slug).Scan(
			&thread.ID,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug)
	}
	return thread, err
}

func (t *ThreadRepository) GetThreadsByForumSlug(slug, since, desc string, limit int) ([]*models.Thread, error) {
	query := `SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created from threads as t
  				LEFT JOIN forum f on t.forum = f.slug
				WHERE f.slug = $1`
	if since != "" && desc == "true" {
		query += " and t.created <= $2"
	} else if since != "" && desc == "false" {
		query += " and t.created >= $2"
	} else if since != "" {
		query += " and t.created >= $2"
	}
	if desc == "true" {
		query += " ORDER BY t.created desc"
	} else if desc == "false" {
		query += " ORDER BY t.created asc"
	} else {
		query += " ORDER BY t.created"
	}
	query += fmt.Sprintf(" LIMIT NULLIF(%d, 0)", limit)
	var rows pgx.Rows
	var err error
	if since != "" {
		rows, err = t.Con.Query(context.Background(), query, slug, since)
	} else {
		rows, err = t.Con.Query(context.Background(), query, slug)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threads := make([]*models.Thread, 0)
	for rows.Next() {
		t := &time.Time{}
		thread := &models.Thread{}
		err = rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			t)
		thread.Created = strfmt.DateTime(t.UTC()).String()
		threads = append(threads, thread)
	}
	return threads, nil
}

func (tr *ThreadRepository) FindThreadByID(threadID int) (*models.Thread, error) {
	query := `SELECT id, title, author, forum, message, votes, slug, created FROM Threads
			  WHERE id = $1`
	thread := &models.Thread{}
	t := &time.Time{}

	err := tr.Con.QueryRow(
		context.Background(),
		query,
		threadID).Scan(
		&thread.ID,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		t)
	thread.Created = strfmt.DateTime(t.UTC()).String()
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (tr *ThreadRepository) CheckThreadID(parentID int) (int, error) {
	query := `SELECT thread FROM Posts WHERE id = $1`
	var threadID int

	err := tr.Con.QueryRow(context.Background(), query, parentID).Scan(&threadID)
	return threadID, err
}

func (tr *ThreadRepository) CreatePost(posts []*postModels.Post) ([]*postModels.Post, error) {
	query := `INSERT INTO Posts (parent, author, message, forum, thread) VALUES `
	for i, post := range posts {
		if i != 0 {
			query += ", "
		}
		query += fmt.Sprintf("(NULLIF(%d, 0), '%s', '%s', '%s', %d)",
			post.Parent,
			post.Author,
			post.Message,
			post.Forum,
			post.Thread)
	}
	query += " RETURNING id, parent, author, message, is_edited, forum, thread, created"

	rows, err := tr.Con.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	newPosts := make([]*postModels.Post, 0)
	var parent sql.NullInt64
	defer rows.Close()
	for rows.Next() {
		t := &time.Time{}
		post := &postModels.Post{}
		err = rows.Scan(
			&post.ID,
			&parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			t)
		if err != nil {
			fmt.Println(err)
		}
		if parent.Valid {
			post.Parent = int(parent.Int64)
		}
		if err != nil {
			return nil, err
		}
		post.Created = strfmt.DateTime(t.UTC()).String()
		newPosts = append(newPosts, post)
	}
	return newPosts, nil
}

func (tr *ThreadRepository) UpdateThreadByID(thread *models.Thread) (*models.Thread, error) {
	query := `UPDATE Threads SET title = $1, message = $2 WHERE id = $3
			  RETURNING id, title, author, forum, message, votes, slug, created`

	t := &time.Time{}
	err := tr.Con.QueryRow(
		context.Background(),
		query,
		thread.Title,
		thread.Message,
		thread.ID).Scan(
		&thread.ID,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		t)
	if err != nil {
		return nil, err
	}
	thread.Created = strfmt.DateTime(t.UTC()).String()

	return thread, nil
}

func (tr *ThreadRepository) UpdateThreadBySlug(thread *models.Thread) (*models.Thread, error) {
	query := `UPDATE Threads SET title = $1, message = $2 WHERE slug = $3
			  RETURNING id, title, author, forum, message, votes, slug, created`

	t := &time.Time{}
	err := tr.Con.QueryRow(
		context.Background(),
		query,
		thread.Title,
		thread.Message,
		thread.Slug).Scan(
		&thread.ID,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		t)
	if err != nil {
		return nil, err
	}
	thread.Created = strfmt.DateTime(t.UTC()).String()

	return thread, nil
}

func (tr *ThreadRepository) GetPosts(limit, threadID int, sort, since string, desc bool) ([]*postModels.Post, error) {
	postID, _ := strconv.Atoi(since)
	var query string

	if sort == "flat" || sort == "" {
		query = FormQueryFlatSort(limit, threadID, sort, since, desc)
	} else if sort == "tree" {
		query = FormQuerySortTree(limit, threadID, postID, sort, since, desc)
	} else if sort == "parent_tree" {
		query = FormQuerySortParentTree(limit, threadID, postID, sort, since, desc)
	}
	if sort != "parent_tree" {
		query += fmt.Sprintf(" LIMIT NULLIF(%d, 0)", limit)
	}
	var rows pgx.Rows
	var err error
	if sort == "tree" {
		rows, err = tr.Con.Query(context.Background(), query, threadID)
	} else if sort == "parent_tree" {
		rows, err = tr.Con.Query(context.Background(), query, threadID)
	} else if since != "" {
		rows, err = tr.Con.Query(context.Background(), query, threadID, postID)
	} else {
		rows, err = tr.Con.Query(context.Background(), query, threadID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*postModels.Post, 0)
	var parent sql.NullInt64
	for rows.Next() {
		t := &time.Time{}
		post := &postModels.Post{}
		err = rows.Scan(&post.ID, &parent, &post.Author, &post.Message,
			&post.IsEdited, &post.Forum, &post.Thread, t)
		post.Created = strfmt.DateTime(t.UTC()).String()
		if parent.Valid {
			post.Parent = int(parent.Int64)
		}
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func FormQueryFlatSort(limit, threadID int, sort, since string, desc bool) string {
	query := `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Posts WHERE thread = $1`
	if since != "" && desc {
		query += " and id < $2"
	} else if since != "" && !desc {
		query += " and id > $2"
	} else if since != "" {
		query += " and id > $2"
	}
	if desc {
		query += " ORDER BY created DESC, id DESC"
	} else if !desc {
		query += " ORDER BY created ASC, id"
	} else {
		query += " ORDER BY created, id"
	}
	return query
}

func FormQuerySortTree(limit, threadID, ID int, sort, since string, desc bool) string {
	query := `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Posts WHERE thread = $1`
	if since != "" && desc {
		query += " and path < "
	} else if since != "" && !desc {
		query += " and path > "
	} else if since != "" {
		query += " and path > "
	}
	if since != "" {
		query += fmt.Sprintf(` (SELECT path FROM Posts WHERE id = %d) `, ID)
	}
	if desc {
		query += " ORDER BY path DESC"
	} else if !desc {
		query += " ORDER BY path ASC, id"
	} else {
		query += " ORDER BY path, id"
	}
	return query
}

func FormQuerySortParentTree(limit, threadID, ID int, sort, since string, desc bool) string {
	subQuery := ""
	query := `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Posts WHERE path[1] IN `
	if since != "" {
		if desc {
			subQuery = `and path[1] < `
		} else {
			subQuery = `and path[1] > `
		}
		subQuery += fmt.Sprintf(`(SELECT path[1] FROM Posts WHERE id = %d)`, ID)
	}
	subQuery = `SELECT id FROM Posts WHERE thread = $1 AND parent is null ` + subQuery
	if desc {
		subQuery += `ORDER BY id DESC`
		subQuery += fmt.Sprintf(` LIMIT NULLIF(%d, 0)`, limit)
		query += fmt.Sprintf(`(%s) ORDER BY path[1] DESC, path, id`, subQuery)
	} else {
		subQuery += `ORDER BY id ASC`
		subQuery += fmt.Sprintf(` LIMIT NULLIF(%d, 0)`, limit)
		query += fmt.Sprintf(`(%s) ORDER BY path, id`, subQuery)
	}
	return query
}

func (tr *ThreadRepository) CreateNewVote(vote *models.Vote) error {
	query := `INSERT INTO Votes (nickname, thread_id, voice)
			  VALUES ($1, $2, $3)`

	_, err := tr.Con.Exec(
		context.Background(),
		query,
		vote.Nickname,
		vote.ThreadID,
		vote.Voice)
	return err
}

func (tr *ThreadRepository) UpdateVote(vote *models.Vote) (int, error) {
	query := `UPDATE Votes SET voice = $1
		      WHERE thread_id = $2 AND nickname = $3 AND voice != $1`

	res, err := tr.Con.Exec(
		context.Background(),
		query,
		vote.Voice,
		vote.ThreadID,
		vote.Nickname)
	return int(res.RowsAffected()), err
}
