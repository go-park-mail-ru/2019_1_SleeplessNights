package user_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"golang.org/x/net/context"
	"math"
)

var LeaderBoardLen = config.GetInt("user_ms.pkg.user_manager.board_len")
var PageLen = uint64(config.GetInt("user_ms.pkg.user_manager.page_len"))

func (us *userManager) GetLeaderBoardPage(ctx context.Context, in *services.PageData) (*services.LeaderBoardPage, error) {
	var page *services.LeaderBoardPage
	for i := in.Since*PageLen - PageLen; i < in.Since*PageLen-1; i++ {
		page.Leaders = append(page.Leaders, profiles[i])
	}
	if len(profiles) < LeaderBoardLen {
		page.PagesCount = uint64(math.Ceil(float64(len(profiles)) / float64(PageLen)))
	}
	return page, nil
}
