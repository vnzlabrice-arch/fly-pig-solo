package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"user-api/internal/logic"
	"user-api/internal/svc"
	"user-api/internal/types"
)

func GetUnpaidOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUnpaidOrdersReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetUnpaidOrdersLogic(r.Context(), svcCtx)
		resp, err := l.GetUnpaidOrders(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
