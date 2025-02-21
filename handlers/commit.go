package handlers

import (
	"strings"

	"github.com/tanayarun/lazydev/services"

	"github.com/gofiber/fiber/v2"
)

func GetCommitHandler(c *fiber.Ctx) error {
	owner := c.Query("owner", "torvalds") // Default repo
	repo := c.Query("repo", "linux")

	commit, err := services.FetchLatestCommit(owner, repo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch commit"})
	}

	return c.JSON(fiber.Map{
		"sha":     commit.Sha,
		"message": strings.Split(commit.Commit.Message, "\n")[0], // Get only the first line
	})

}
