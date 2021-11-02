package httpd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mvt-demo/internal/logger"
	"mvt-demo/internal/model"
	"mvt-demo/internal/params"
	"mvt-demo/internal/utils/response"
	"mvt-demo/internal/utils/xyz2lonlat"
	"net/http"
)

// @Summary mvt
// @Description  mvt
// @version 1.0
// @tags mvt
// @Param z path int true "z"
// @Param x path int true "x"
// @Param y path int true "y"
// @Success 200 {object} string "{"code":200,"data": "","msg":"success"}"
// @Failure 400 {string} json "{"code":400,"data":{},"msg":"bind query err/params error"}"
// @Router /mvt/{z}/{x}/{y} [get]
func MVTGet(c *gin.Context) {
	utilGin := response.Gin{Ctx: c}
	p := params.MVTGet{}
	if err := c.BindUri(&p); err != nil {
		logger.Logger.Error("bind uri err",
			zap.Error(err))
		utilGin.Response(400, "bind uri err", nil)
		return
	}
	xyMin := xyz2lonlat.XYZ2lonlat(p.X, p.Y, p.Z)
	xyMax := xyz2lonlat.XYZ2lonlat(p.X+1, p.Y+1, p.Z)

	var mvt []byte
//	sql4326 := fmt.Sprintf(`select ST_AsMVT(P, '%v', 4096, 'geom') as "mvt" from (
//select ST_AsMVTGeom(ST_Transform(geom, 4326),ST_MakeEnvelope(%v,%v,%v,%v, 4326),
//4096, 64, TRUE) geom FROM adcode
//where ST_Transform(geom, 4326) && ST_MakeEnvelope(%v,%v,%v,%v, 4326)
//) AS P;`, "adcode", xyMin[0], xyMin[1], xyMax[0], xyMax[1], xyMin[0], xyMin[1], xyMax[0], xyMax[1])
//	sql := sql4326

//	sql3857 := fmt.Sprintf(`select ST_AsMVT(P, '%v', 4096, 'geom') as "mvt" from (
//select fid, ST_AsMVTGeom(ST_Transform(geom, 3857),ST_Transform(ST_MakeEnvelope(%v,%v,%v,%v, 4326), 3857),
//4096, 64, TRUE) geom FROM adcode
//where ST_Transform(geom, 4326) && ST_MakeEnvelope(%v,%v,%v,%v, 4326)
//) AS P;`, "adcode", xyMin[0], xyMin[1], xyMax[0], xyMax[1], xyMin[0], xyMin[1], xyMax[0], xyMax[1])
//	sql := sql3857

	sqlMany := fmt.Sprintf(`select ((select ST_AsMVT(P, '%v', 4096, 'geom') as "mvt" from (
select fid, ST_AsMVTGeom(ST_Transform(geom, 3857),ST_Transform(ST_MakeEnvelope(%v,%v,%v,%v, 4326), 3857),
4096, 64, TRUE) geom FROM adcode
where ST_Transform(geom, 4326) && ST_MakeEnvelope(%v,%v,%v,%v, 4326)
) AS P)
|| (select ST_AsMVT(P, '%v', 4096, 'geom') as "mvt" from (
select fid, ST_AsMVTGeom(ST_Transform(geom, 3857),ST_Transform(ST_MakeEnvelope(%v,%v,%v,%v, 4326), 3857),
4096, 64, TRUE) geom FROM adcode_many
where ST_Transform(geom, 4326) && ST_MakeEnvelope(%v,%v,%v,%v, 4326)
) AS P)) as mvt;`, "adcode", xyMin[0], xyMin[1], xyMax[0], xyMax[1], xyMin[0], xyMin[1], xyMax[0], xyMax[1],
		"adcode_many", xyMin[0], xyMin[1], xyMax[0], xyMax[1], xyMin[0], xyMin[1], xyMax[0], xyMax[1])
	sql := sqlMany

	row := model.DB.Debug().Raw(sql).Row()
	row.Scan(&mvt)
	c.Header("Content-Length", fmt.Sprintf("%d", len(mvt)))
	c.Data(http.StatusOK, "application/x-protobuf", mvt)
	return
}
