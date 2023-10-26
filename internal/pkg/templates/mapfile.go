package templates

const MapfileTemplate = `#  Tile Index
LAYER
 DEBUG ON
 STATUS ON
 NAME "{{.Product}}_time_idx"
 TYPE POLYGON
 DATA "geom from (SELECT * FROM raster_geoms WHERE product = '{{.Product}}') as subquery USING unique gid USING srid=4326"
 METADATA
  "wms_title" "TIME INDEX"
  "wms_srs"   "EPSG:4326"
  "wms_extent" "-60 50 20 90"
  INCLUDE "times_{{.Product}}.map"
  "wms_timeextent "{{.AllTimesString}},{{.StartRange}}/{{.EndRange}}
	"wms_timedefault" "{{.DefaultTimeString}}"
  "wms_timeitem" "datetime" #column in postgis table of type timestamp
  "wms_enable_request" "*"
 END
 CONNECTION "host=<PSQL_HOST> dbname=<PSQL_DB> user=<PSQL_USER> password=<PSQL_PASS>"
 CONNECTIONTYPE postgis
END

# raster layer
LAYER
 NAME "{{.Product}}"
 TYPE RASTER
 STATUS ON
 DEBUG ON
 PROJECTION
   "init=epsg:4326"
 END
 METADATA
  "wms_title" "{{.Product}}"
  "wms_srs"   "EPSG:4326"
  INCLUDE "times_{{.Product}}.map"
  "wms_timeextent "{{.AllTimesString}},{{.StartRange}}/{{.EndRange}}
  "wms_timedefault" "{{.DefaultTimeString}}"
  "wms_extent" "-60 50 20 90"
  "wms_timeitem" "datetime" #datetime is a column in postgis table of type timestamp
  "wms_enable_request" "*"
 END
 #OFFSITE 0 0 0
 TILEITEM "location" #filepath is a column in postgis table with varchar of the filepath to each image
 TILEINDEX "{{.Product}}_time_idx"
 TILESRS "src_srs"
END`
