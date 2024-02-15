package tests

import (
	"bytes"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/xuxant/kbuild/cmd/kbuild"
	"testing"
)

type TestSuite struct {
	suite.Suite
	ctx context.Context
}

func (suite *TestSuite) SetupTest() {
	log.Info().Msg("Setting up tests...")
}

func (suite *TestSuite) TearDownTest() {
	log.Info().Msg("Tearing down tests...")
}

func TestRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) UnitTest() {
	t := suite.T()
	log.Info().Msg("Running unit tests....")
	fmt.Println("yuck")
	var output bytes.Buffer
	suite.Run("test kbuild", func() {
		fmt.Println("yuck")
		command := kbuild.NewBuildCommand()
		command.SetOut(&output)

		command.Execute()

		expectedOutput := "Remote build with Kaniko\n"
		require.NotNil(t, output.String())
		require.Contains(t, expectedOutput, output.String())
	})
}
