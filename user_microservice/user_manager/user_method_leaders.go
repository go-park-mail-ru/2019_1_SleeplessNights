package user_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"golang.org/x/net/context"
	"math"
)

var PageLen = uint64(config.GetInt("user_ms.pkg.user_manager.page_len"))

func (us *userManager) GetLeaderBoardPage(ctx context.Context, in *services.PageData) (*services.LeaderBoardPage, error) {
	var page services.LeaderBoardPage
	for i, p := range profiles {
		if uint64(i) < in.Since*10-10 {
			continue
		}
		if uint64(i) == in.Since*10 {
			break
		}
		page.Leaders = append(page.Leaders, p)
	}
	l := uint64(len(profiles))
	if l < LeaderBoardLen {
		page.PagesCount = uint64(math.Ceil(float64(l) / float64(PageLen)))
	}
	return &page, nil
}
