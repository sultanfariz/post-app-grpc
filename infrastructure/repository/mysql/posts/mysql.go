package posts

import (
	"context"
	"time"

	model "github.com/sultanfariz/simple-grpc/domain/posts"

	"gorm.io/gorm"
)

type PostsRepository struct {
	DBConnection *gorm.DB
}

func NewPostsRepository(db *gorm.DB) *PostsRepository {
	return &PostsRepository{
		DBConnection: db,
	}
}

func (r *PostsRepository) GetAll(ctx context.Context) ([]*model.Post, error) {
	data := []*model.Post{}
	if err := r.DBConnection.Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *PostsRepository) GetById(ctx context.Context, id int) (*model.Post, error) {
	data := model.Post{}
	if err := r.DBConnection.Where("id = ?", id).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *PostsRepository) Insert(ctx context.Context, post *model.Post) (*model.Post, error) {
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	if err := r.DBConnection.Create(&post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostsRepository) Delete(ctx context.Context, id int) error {
	if err := r.DBConnection.Where("id = ?", id).Delete(&model.Post{}).Error; err != nil {
		return err
	}

	return nil
}
