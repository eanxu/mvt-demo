该demo实现mapbox矢量瓦片服务,结合postgis, gin, gorm

文件说明:   
adcode.zip 为示例数据, 需导入postgresql数据库中,并存储两份 命名: adcode与adcode_many(该表可删除部分数据,用于显示单服务多图层)
meta.json 为根据mapbox提供的工具形成的服务元数据样例
mvt.html 为瓦片加载前端代码
