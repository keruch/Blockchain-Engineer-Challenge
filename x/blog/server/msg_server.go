package server

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/regen-network/bec/x/blog"
)

var _ blog.MsgServer = serverImpl{}

func (s serverImpl) CreatePost(goCtx context.Context, request *blog.MsgCreatePost) (*blog.MsgCreatePostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := prefix.NewStore(ctx.KVStore(s.storeKey), blog.KeyPrefix(blog.PostKey))

	key := []byte(request.Slug)
	if store.Has(key) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "duplicate slug %s found", request.Slug)
	}

	post := blog.Post{
		Author: request.Author,
		Slug:   request.Slug,
		Title:  request.Title,
		Body:   request.Body,
	}

	bz, err := s.cdc.Marshal(&post)
	if err != nil {
		return nil, err
	}

	store.Set(key, bz)

	return &blog.MsgCreatePostResponse{}, nil
}

func (s serverImpl) CreateComment(goCtx context.Context, request *blog.MsgCreateComment) (*blog.MsgCreateCommentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	postStore := prefix.NewStore(ctx.KVStore(s.storeKey), blog.KeyPrefix(blog.PostKey))

	// check if the specified post exists
	postKey := []byte(request.PostSlug)
	if !postStore.Has(postKey) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "post with slug %s not found", request.PostSlug)
	}

	comment := blog.Comment{
		PostSlug: request.PostSlug,
		Author:   request.Author,
		Body:     request.Body,
	}

	bz, err := s.cdc.Marshal(&comment)
	if err != nil {
		return nil, err
	}

	commentKey := generateCommentKey(request.PostSlug)

	commentStore := prefix.NewStore(ctx.KVStore(s.storeKey), blog.KeyPrefix(blog.CommentKey))

	commentStore.Set(commentKey, bz)

	return &blog.MsgCreateCommentResponse{}, nil
}
