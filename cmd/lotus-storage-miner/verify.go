package main

import (
	_ "net/http/pprof"
	"time"

	"github.com/filecoin-project/lotus/api"
	lcli "github.com/filecoin-project/lotus/cli"
	"github.com/filecoin-project/lotus/node"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/filecoin-project/lotus/storage"
	"github.com/urfave/cli/v2"
)

var verifyCmd = &cli.Command{
	Name:  "verify",
	Usage: "verify sectors",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "di",
			Usage: "deadline index",
			Value: -1,
		},
	},
	Action: func(cctx *cli.Context) error {
		storage.ForceChangeDI = cctx.Int("di")

		nodeApi, ncloser, err := lcli.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer ncloser()
		ctx := lcli.DaemonContext(cctx)

		minerRepoPath := cctx.String(FlagMinerRepo)
		r, err := repo.NewFS(minerRepoPath)
		if err != nil {
			return err
		}

		var minerapi api.StorageMiner
		_, err = node.New(ctx,
			node.StorageMiner(&minerapi),
			node.Online(),
			node.Repo(r),
			node.Override(new(api.FullNode), nodeApi),
		)

		// Bootstrap with full node
		remoteAddrs, err := nodeApi.NetAddrsListen(ctx)
		if err != nil {
			return err
		}

		if err := minerapi.NetConnect(ctx, remoteAddrs); err != nil {
			return err
		}

		<-time.After(10 * time.Hour)
		return err
	},
}
