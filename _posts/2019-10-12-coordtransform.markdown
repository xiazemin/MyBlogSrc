---
title: coordtransform
layout: post
category: algorithm
author: 夏泽民
---
https://github.com/qichengzx/coordtransform
https://github.com/wandergis/coordtransform

https://gist.github.com/jp1017/71bd0976287ce163c11a7cb963b04dd8
当前互联网地图的坐标系现状
地球坐标 (WGS84)
国际标准，从 GPS 设备中取出的数据的坐标系
国际地图提供商使用的坐标系
火星坐标 (GCJ-02)也叫国测局坐标系
中国标准，从国行移动设备中定位获取的坐标数据使用这个坐标系
国家规定： 国内出版的各种地图系统（包括电子形式），必须至少采用GCJ-02对地理位置进行首次加密。
百度坐标 (BD-09)
百度标准，百度 SDK，百度地图，Geocoding 使用
(本来就乱了，百度又在火星坐标上来个二次加密)
开发过程需要注意的事
从设备获取经纬度（GPS）坐标

  	如果使用的是百度sdk那么可以获得百度坐标（bd09）或者火星坐标（GCJ02),默认是bd09
  	如果使用的是ios的原生定位库，那么获得的坐标是WGS84
  	如果使用的是高德sdk,那么获取的坐标是GCJ02
互联网在线地图使用的坐标系

  火星坐标系：
  		iOS 地图（其实是高德）
  		Google国内地图（.cn域名下）
  		搜搜、阿里云、高德地图、腾讯
  百度坐标系：
  		当然只有百度地图
  WGS84坐标系：
  		国际标准，谷歌国外地图、osm地图等国外的地图一般都是这个
<!-- more -->
1.地球坐标 (WGS84)
1.WGS84是现行的国际标准，例如从iphone中 GPS 设备中取出的数据就是WGS84坐标系。
2.国际地图提供商使用的坐标系。

2.国家测汇局坐标系或者火星坐标系(GCJ-02)
1.中国标准，从国行移动设备中定位获取的坐标数据使用这个坐标系
2.国家规定： 国内出版的各种地图系统（包括电子形式），必须至少采用GCJ-02对地理位置进行首次加密。

3.百度坐标 (BD-09)
1.国内百度的标准，百度 SDK，百度地图，百度相关的产品使用

4.各个坐标系之间的转换逻辑
public class GpsUtils {

    private static double x_pi = 3.14159265358979324 * 3000.0 / 180.0;
    // π
    private static double pi = 3.1415926535897932384626;
    // 长半轴
    private static double a = 6378245.0;
    // 扁率
    private static double ee = 0.00669342162296594323;

    public static boolean out_of_china(double lon, double lat) {
        if(lon < 72.004 || lon > 137.8347) {
            return true;
        } else if(lat < 0.8293 || lat > 55.8271) {
            return true;
        }
        return false;
    }

    public static double transformlat(double lon, double lat) {
        double ret = -100.0 + 2.0 * lon + 3.0 * lat + 0.2 * lat * lat + 0.1 * lon * lat + 0.2 * Math.sqrt(Math.abs(lon));
        ret += (20.0 * Math.sin(6.0 * lon * pi) + 20.0 * Math.sin(2.0 * lon * pi)) * 2.0 / 3.0;
        ret += (20.0 * Math.sin(lat * pi) + 40.0 * Math.sin(lat / 3.0 * pi)) * 2.0 / 3.0;
        ret += (160.0 * Math.sin(lat / 12.0 * pi) + 320 * Math.sin(lat * pi / 30.0)) * 2.0 / 3.0;
        return ret;
    }

    public static double transformlng(double lon, double lat) {
        double ret = 300.0 + lon + 2.0 * lat + 0.1 * lon * lon + 0.1 * lon * lat + 0.1 * Math.sqrt(Math.abs(lon));
        ret += (20.0 * Math.sin(6.0 * lon * pi) + 20.0 * Math.sin(2.0 * lon * pi)) * 2.0 / 3.0;
        ret += (20.0 * Math.sin(lon * pi) + 40.0 * Math.sin(lon / 3.0 * pi)) * 2.0 / 3.0;
        ret += (150.0 * Math.sin(lon / 12.0 * pi) + 300.0 * Math.sin(lon / 30.0 * pi)) * 2.0 / 3.0;
        return ret;
    }

    /**
     * WGS84转GCJ02(火星坐标系)
     *
     * @param wgs_lon WGS84坐标系的经度
     * @param wgs_lat WGS84坐标系的纬度
     * @return 火星坐标数组
     */
    public static double[] wgs84togcj02(double wgs_lon, double wgs_lat) {
        if (out_of_china(wgs_lon, wgs_lat)) {
            return new double[] { wgs_lon, wgs_lat };
        }
        double dlat = transformlat(wgs_lon - 105.0, wgs_lat - 35.0);
        double dlng = transformlng(wgs_lon - 105.0, wgs_lat - 35.0);
        double radlat = wgs_lat / 180.0 * pi;
        double magic = Math.sin(radlat);
        magic = 1 - ee * magic * magic;
        double sqrtmagic = Math.sqrt(magic);
        dlat = (dlat * 180.0) / ((a * (1 - ee)) / (magic * sqrtmagic) * pi);
        dlng = (dlng * 180.0) / (a / sqrtmagic * Math.cos(radlat) * pi);
        double mglat = wgs_lat + dlat;
        double mglng = wgs_lon + dlng;
        return new double[] { mglng, mglat };
    }

    /**
     * GCJ02(火星坐标系)转GPS84
     *
     * @param gcj_lon 火星坐标系的经度
     * @param gcj_lat 火星坐标系纬度
     * @return WGS84坐标数组
     */
    public static double[] gcj02towgs84(double gcj_lon, double gcj_lat) {
        if (out_of_china(gcj_lon, gcj_lat)) {
            return new double[] { gcj_lon, gcj_lat };
        }
        double dlat = transformlat(gcj_lon - 105.0, gcj_lat - 35.0);
        double dlng = transformlng(gcj_lon - 105.0, gcj_lat - 35.0);
        double radlat = gcj_lat / 180.0 * pi;
        double magic = Math.sin(radlat);
        magic = 1 - ee * magic * magic;
        double sqrtmagic = Math.sqrt(magic);
        dlat = (dlat * 180.0) / ((a * (1 - ee)) / (magic * sqrtmagic) * pi);
        dlng = (dlng * 180.0) / (a / sqrtmagic * Math.cos(radlat) * pi);
        double mglat = gcj_lat + dlat;
        double mglng = gcj_lon + dlng;
        return new double[] { gcj_lon * 2 - mglng, gcj_lat * 2 - mglat };
    }


    /**
     * 火星坐标系(GCJ-02)转百度坐标系(BD-09)
     *
     * 谷歌、高德——>百度
     * @param gcj_lon 火星坐标经度
     * @param gcj_lat 火星坐标纬度
     * @return 百度坐标数组
     */
    public static double[] gcj02tobd09(double gcj_lon, double gcj_lat) {
        double z = Math.sqrt(gcj_lon * gcj_lon + gcj_lat * gcj_lat) + 0.00002 * Math.sin(gcj_lat * x_pi);
        double theta = Math.atan2(gcj_lat, gcj_lon) + 0.000003 * Math.cos(gcj_lon * x_pi);
        double bd_lng = z * Math.cos(theta) + 0.0065;
        double bd_lat = z * Math.sin(theta) + 0.006;
        return new double[] { bd_lng, bd_lat };
    }

    /**
     * 百度坐标系(BD-09)转火星坐标系(GCJ-02)
     *
     * 百度——>谷歌、高德
     * @param bd_lon 百度坐标纬度
     * @param bd_lat 百度坐标经度
     * @return 火星坐标数组
     */
    public static double[] bd09togcj02(double bd_lon, double bd_lat) {
        double x = bd_lon - 0.0065;
        double y = bd_lat - 0.006;
        double z = Math.sqrt(x * x + y * y) - 0.00002 * Math.sin(y * x_pi);
        double theta = Math.atan2(y, x) - 0.000003 * Math.cos(x * x_pi);
        double gg_lng = z * Math.cos(theta);
        double gg_lat = z * Math.sin(theta);
        return new double[] { gg_lng, gg_lat };
    }

    /**
     * WGS坐标转百度坐标系(BD-09)
     *
     * @param wgs_lng WGS84坐标系的经度
     * @param wgs_lat WGS84坐标系的纬度
     * @return 百度坐标数组
     */
    public static double[] wgs84tobd09(double wgs_lng, double wgs_lat) {
        double[] gcj = wgs84togcj02(wgs_lng, wgs_lat);
        double[] bd09 = gcj02tobd09(gcj[0], gcj[1]);
        return bd09;
    }

    /**
     * 百度坐标系(BD-09)转WGS坐标
     *
     * @param bd_lng 百度坐标纬度
     * @param bd_lat 百度坐标经度
     * @return WGS84坐标数组
     */
    public static double[] bd09towgs84(double bd_lng, double bd_lat) {
        double[] gcj = bd09togcj02(bd_lng, bd_lat);
        double[] wgs84 = gcj02towgs84(gcj[0], gcj[1]);
        return wgs84;
    }
}
