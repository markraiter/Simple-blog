package models

import "testing"

var (
	testTitile = "This is a test post title"
	testEmail  = "user@example.org"
	testBody   = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum euismod maximus nulla id tincidunt. Aenean vel mi leo. Quisque nec lacinia ex. Suspendisse lobortis mi et neque semper pretium. Proin pellentesque tortor at elit elementum semper. Aenean dictum justo ac urna iaculis luctus. Pellentesque consectetur neque lectus, id tempor ipsum."
)

func TestUser(t *testing.T) *User {
	t.Helper()

	return &User{
		Email:    "user@example.org",
		Password: "password",
	}
}

func TestUpdatePostInput(t *testing.T) *UpdatePostInput {
	t.Helper()

	return &UpdatePostInput{
		Title: &testTitile,
		Body:  &testBody,
	}
}

func TestUpdateCommentInput(t *testing.T) *UpdateCommentInput {
	t.Helper()

	return &UpdateCommentInput{
		Email: &testEmail,
		Body:  &testBody,
	}
}
