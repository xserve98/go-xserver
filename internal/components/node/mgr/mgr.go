package nodemgr

import (
	"os"
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
	"github.com/fananchong/go-xserver/internal/components/misc"
	nodecommon "github.com/fananchong/go-xserver/internal/components/node/common"
	"github.com/fananchong/go-xserver/internal/db"
)

// Mgr : 管理节点
type Mgr struct {
	*nodecommon.Node
	ctx *common.Context
}

// NewMgr : 管理节点实现类的构造函数
func NewMgr(ctx *common.Context) *Mgr {
	mgr := &Mgr{
		ctx:  ctx,
		Node: nodecommon.NewNode(ctx, common.Mgr),
	}
	pluginType := misc.GetPluginType(mgr.ctx.Ctx)
	if pluginType == common.Mgr {
		mgr.init()
		mgr.ctx.Node = mgr
	}
	return mgr
}

// Start : 实例化组件
func (mgr *Mgr) Start() bool {
	pluginType := misc.GetPluginType(mgr.ctx.Ctx)
	if pluginType == common.Mgr {
		if mgr.Node.Start() == false {
			mgr.ctx.Log.Errorln("Mgr Server node failed to start")
			os.Exit(1)
		}
	}
	return true
}

// Close : 关闭组件
func (mgr *Mgr) Close() {
	if mgr.Node != nil {
		mgr.Node.Close()
		mgr.Node = nil
	}
}

func (mgr *Mgr) init() bool {
	// register ticker
	registerTicker := utils.NewTickerHelper("REGISTER", mgr.ctx, 1*time.Second, mgr.register)
	return mgr.Node.Init(Session{}, []utils.IComponent{registerTicker})
}

func (mgr *Mgr) register() {
	data := db.NewMgrServer(mgr.ctx.Config.DbMgr.Name, 0)
	data.SetAddr(utils.GetIPInner(mgr.ctx))
	data.SetPort(utils.GetIntranetListenPort(mgr.ctx))
	if err := data.Save(); err != nil {
		mgr.ctx.Log.Errorln(err)
	}
}
