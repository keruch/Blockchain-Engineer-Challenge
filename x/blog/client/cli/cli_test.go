package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/regen-network/bec/testutil/network"
	"github.com/regen-network/bec/x/blog"
	"github.com/regen-network/bec/x/blog/client/cli"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := network.DefaultConfig()
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestCreatePost() {
	val0 := s.network.Validators[0]

	testCases := []struct {
		name      string
		args      []string
		expErr    bool
		expErrMsg string
	}{
		{
			"no author",
			[]string{"", "foo", "bar", "baz"},
			true, "no author",
		},
		{
			"no slug",
			[]string{val0.Address.String(), "", "bar", "baz"},
			true, "no slug",
		},
		{
			"no title",
			[]string{val0.Address.String(), "foo", "", "baz"},
			true, "no title",
		},
		{
			"no body",
			[]string{val0.Address.String(), "foo", "bar", ""},
			true, "no body",
		},
		{
			"valid request",
			[]string{val0.Address.String(), "foo", "bar", "baz"},
			false, "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := cli.CmdCreatePost()
			args := append([]string{
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			}, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(val0.ClientCtx, cmd, args)

			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				var txRes sdk.TxResponse
				err := val0.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &txRes)
				s.Require().NoError(err)
				s.Require().Equal(uint32(0), txRes.Code)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCreateComment() {
	val0 := s.network.Validators[0]

	// create a dummy post
	const postSlug = "post_slug"
	argsPost := []string{
		val0.Address.String(), postSlug, "title", "body",
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	cmdPost := cli.CmdCreatePost()

	outPost, errPost := clitestutil.ExecTestCLICmd(val0.ClientCtx, cmdPost, argsPost)
	s.Require().NoError(errPost)
	var txResPost sdk.TxResponse
	errPost = val0.ClientCtx.Codec.UnmarshalJSON(outPost.Bytes(), &txResPost)
	s.Require().NoError(errPost)
	s.Require().Equal(uint32(0), txResPost.Code)

	testCases := []struct {
		name      string
		args      []string
		expErr    bool
		expErrMsg string
	}{
		{
			"no post_slug",
			[]string{val0.Address.String(), "", "body"},
			true, "no post_slug",
		},
		{
			"no author",
			[]string{"", postSlug, "body"},
			true, "no author",
		},
		{
			"no body",
			[]string{val0.Address.String(), postSlug, ""},
			true, "no body",
		},
		{
			"valid request",
			[]string{val0.Address.String(), postSlug, "body"},
			false, "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := cli.CmdCreateComment()
			args := append([]string{
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			}, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(val0.ClientCtx, cmd, args)

			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				var txRes sdk.TxResponse
				err = val0.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &txRes)
				s.Require().NoError(err)
				s.Require().Equal(uint32(0), txRes.Code)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestUnknownCommentPostSlug() {
	val0 := s.network.Validators[0]

	// create one dummy comment with unknown post_slug
	const postSlug = "unknown_slug"
	argsComment := []string{
		val0.Address.String(), postSlug, "body",
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	cmdComment := cli.CmdCreateComment()

	out, err := clitestutil.ExecTestCLICmd(val0.ClientCtx, cmdComment, argsComment)
	s.Require().NoError(err)
	var txRes sdk.TxResponse
	err = val0.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &txRes)
	s.Require().NoError(err)
	s.Require().NotEqual(uint32(0), txRes.Code)
}

func (s *IntegrationTestSuite) TestDuplicateSlugs() {
	val0 := s.network.Validators[0]
	args := []string{
		val0.Address.String(), "alice", "bob", "charlie",
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	cmd := cli.CmdCreatePost()

	// Send the first time.
	out, err := clitestutil.ExecTestCLICmd(val0.ClientCtx, cmd, args)
	s.Require().NoError(err)
	var txRes sdk.TxResponse
	err = val0.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &txRes)
	s.Require().NoError(err)
	s.Require().Equal(uint32(0), txRes.Code)

	// Send the second time.
	out, err = clitestutil.ExecTestCLICmd(val0.ClientCtx, cmd, args)
	s.Require().NoError(err)
	err = val0.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &txRes)
	s.Require().NoError(err)
	s.Require().NotEqual(uint32(0), txRes.Code)
}

func (s *IntegrationTestSuite) TestAllPosts() {
	val0 := s.network.Validators[0]

	// Create two dummy posts.
	cmd := cli.CmdCreatePost()
	defaultArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	out, err := clitestutil.ExecTestCLICmd(val0.ClientCtx, cmd, append(defaultArgs, val0.Address.String(), "foo1", "bar1", "baz1"))
	s.Require().NoError(err)
	var txRes1 sdk.TxResponse
	err = val0.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &txRes1)
	s.Require().NoError(err)
	s.Require().Equal(uint32(0), txRes1.Code)
	s.Require().NoError(s.network.WaitForNextBlock())

	out, err = clitestutil.ExecTestCLICmd(val0.ClientCtx, cmd, append(defaultArgs, val0.Address.String(), "foo2", "bar2", "baz1"))
	s.Require().NoError(err)
	var txRes2 sdk.TxResponse
	err = val0.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &txRes2)
	s.Require().NoError(err)
	s.Require().Equal(uint32(0), txRes2.Code)
	s.Require().NoError(s.network.WaitForNextBlock())

	testCases := []struct {
		name        string
		args        []string
		expErr      bool
		expErrMsg   string
		expNumPosts int
	}{
		{
			"no pagination",
			[]string{"-o=json"},
			false, "", 2,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := cli.CmdAllPosts()
			out, err := clitestutil.ExecTestCLICmd(val0.ClientCtx, cmd, tc.args)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				var res = blog.QueryAllPostsResponse{}
				err := val0.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res)
				s.Require().NoError(err)
				s.Require().Equal(tc.expNumPosts, len(res.Posts))
			}
		})
	}
}

func (s *IntegrationTestSuite) TestAllComments() {
	val0 := s.network.Validators[0]

	// create a dummy post
	const postSlug = "post_slug"
	argsPost := []string{
		val0.Address.String(), postSlug, "title", "body",
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	cmdPost := cli.CmdCreatePost()

	outPost, errPost := clitestutil.ExecTestCLICmd(val0.ClientCtx, cmdPost, argsPost)
	s.Require().NoError(errPost)
	var txResPost sdk.TxResponse
	errPost = val0.ClientCtx.Codec.UnmarshalJSON(outPost.Bytes(), &txResPost)
	s.Require().NoError(errPost)
	s.Require().Equal(uint32(0), txResPost.Code)

	// create two dummy comments
	// 1
	argsComment := []string{
		val0.Address.String(), postSlug, "body1",
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	cmdComment := cli.CmdCreateComment()

	out, err := clitestutil.ExecTestCLICmd(val0.ClientCtx, cmdComment, argsComment)
	s.Require().NoError(err)
	var txRes sdk.TxResponse
	err = val0.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &txRes)
	s.Require().NoError(err)
	s.Require().Equal(uint32(0), txRes.Code)

	// 2
	argsComment = []string{
		val0.Address.String(), postSlug, "body2",
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	cmdComment = cli.CmdCreateComment()

	out, err = clitestutil.ExecTestCLICmd(val0.ClientCtx, cmdComment, argsComment)
	s.Require().NoError(err)
	err = val0.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &txRes)
	s.Require().NoError(err)
	s.Require().Equal(uint32(0), txResPost.Code)

	testCases := []struct {
		name           string
		args           []string
		expErr         bool
		expErrMsg      string
		expNumComments int
	}{
		{
			"no pagination",
			[]string{postSlug, "-o=json"},
			false, "", 2,
		},
		{
			"unknown post_slug",
			[]string{"unknown_slug", "-o=json"},
			false, "", 0,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := cli.CmdAllComments()
			out, err = clitestutil.ExecTestCLICmd(val0.ClientCtx, cmd, tc.args)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				var res = blog.QueryAllCommentsResponse{}
				err = val0.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res)
				s.Require().NoError(err)
				s.Require().Equal(tc.expNumComments, len(res.Comments))
			}
		})
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
