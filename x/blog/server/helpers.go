package server

import (
	"strings"

	"github.com/oklog/ulid/v2"

	"github.com/regen-network/bec/x/blog"
)

// generateCommentKey generates a unique id for the post. It uses ULID, which stands for
// Universally Unique Lexicographically Sortable Identifier (https://github.com/ulid/spec).
// This allows us to sort comments by the creation date, if needed.
// The output format is: post_slug-<ULID>.
func generateCommentKey(postSlug string) []byte {
	return []byte(postSlug + "-" + ulid.Make().String())
}

// commentRelatesToPost checks whether a comment by the given key relates to the given post by its slug.
// This logic is moved to a different method in order to make working with comment keys safer.
func commentRelatesToPost(commentKey []byte, postSlug string) bool {
	return strings.HasPrefix(string(commentKey), blog.CommentKey+postSlug)
}
