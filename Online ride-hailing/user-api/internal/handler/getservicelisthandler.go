package handler

import (
	"net/http"

	"user-api/internal/logic"
	"user-api/internal/svc"
	"user-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetServiceListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ServiceListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetServiceListLogic(r.Context(), svcCtx)
		resp, err := l.GetServiceList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
