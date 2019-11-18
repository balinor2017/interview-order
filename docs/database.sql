-- MySQL dump 10.13  Distrib 5.7.27, for Linux (x86_64)
--
-- Host: hk.imgn.to    Database: orders
-- ------------------------------------------------------
-- Server version	5.7.27-0ubuntu0.18.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `rs_order`
--

DROP TABLE IF EXISTS `rs_order`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rs_order` (
  `order_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `pickup_time` datetime NOT NULL,
  `customer_id` bigint(20) NOT NULL DEFAULT '0',
  `driver_id` bigint(20) NOT NULL DEFAULT '0',
  `pickup_longitude` double NOT NULL,
  `pickup_latitude` double NOT NULL,
  `dropoff_longitude` double NOT NULL,
  `dropoff_latitude` double NOT NULL,
  `status` char(8) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`order_id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rs_order`
--

LOCK TABLES `rs_order` WRITE;
/*!40000 ALTER TABLE `rs_order` DISABLE KEYS */;
INSERT INTO `rs_order` VALUES (1,'2019-08-20 15:04:05',1001,0,40.023316,116.459461,42.023316,116.459461,'0','2019-08-21 13:52:08','2019-08-21 13:52:08'),(2,'2019-08-20 15:04:05',1001,0,40.023316,116.459461,42.023316,116.459461,'0','2019-08-21 14:28:23','2019-08-21 14:28:23'),(3,'2019-08-20 15:04:05',111,0,40.023316,116.459461,42.023316,116.459461,'0','2019-08-21 15:38:56','2019-08-21 15:38:56'),(4,'2019-08-20 15:04:05',1001,0,40.023316,116.459461,42.023316,116.459461,'0','2019-08-21 15:39:18','2019-08-21 15:39:18'),(5,'2019-08-20 15:04:05',1001,0,36.069554,120.429545,36.070039,120.431989,'0','2019-08-21 15:47:07','2019-08-21 15:47:07'),(6,'2019-08-20 15:04:05',1001,2001,36.069554,120.429545,36.080039,120.431989,'4','2019-08-21 15:54:47','2019-08-22 11:50:36'),(7,'2019-08-20 15:04:05',1001,0,36.069554,120.429545,36.080039,120.431989,'0','2019-08-21 16:18:31','2019-08-21 16:18:31'),(8,'2019-08-20 15:04:05',1001,0,36.069554,120.429545,36.080039,120.431989,'0','2019-08-21 16:33:24','2019-08-21 16:33:24'),(9,'2019-08-20 15:04:05',1001,0,36.069554,120.429545,36.080039,120.431989,'0','2019-08-21 17:09:52','2019-08-21 17:09:52'),(10,'2019-08-20 15:04:05',1001,2001,36.069554,120.429545,36.080039,120.431989,'0','2019-08-22 07:46:37','2019-08-22 08:00:57'),(11,'2019-08-20 15:04:05',1001,0,36.069554,120.429545,36.080039,120.431989,'0','2019-08-22 09:45:42','2019-08-22 09:45:42'),(12,'2019-08-22 15:04:05',1001,2001,36.06874,120.428701,36.106803,120.362873,'4','2019-08-22 12:03:05','2019-08-23 14:34:17'),(13,'2019-08-22 15:04:05',1001,0,36.06874,120.428701,36.106803,120.362873,'0','2019-08-23 11:38:47','2019-08-23 11:38:47'),(14,'2019-08-22 15:04:05',1001,0,36.06874,120.428701,36.106803,120.362873,'0','2019-08-23 11:39:14','2019-08-23 11:39:14'),(15,'2019-08-22 15:04:05',1002,2001,36.06874,120.428701,36.106803,120.362873,'4','2019-08-26 00:29:21','2019-08-26 02:22:45'),(16,'2019-08-22 15:04:05',1001,2001,36.06874,120.428701,36.106803,120.362873,'4','2019-08-26 06:32:54','2019-08-26 06:49:25');
/*!40000 ALTER TABLE `rs_order` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rs_trip`
--

DROP TABLE IF EXISTS `rs_trip`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rs_trip` (
  `trip_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `customer_id` bigint(20) NOT NULL,
  `order_id` bigint(20) NOT NULL,
  `driver_id` bigint(20) NOT NULL,
  `start_time` datetime NOT NULL,
  `end_time` datetime NOT NULL,
  `vacant_distance` int(11) NOT NULL,
  `engaged_distance` int(11) NOT NULL,
  `pickup_longitude` double NOT NULL,
  `pickup_latitude` double NOT NULL,
  `dropoff_longitude` double NOT NULL,
  `dropoff_latitude` double NOT NULL,
  `Status` int(11) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`trip_id`,`order_id`),
  UNIQUE KEY `order_id` (`order_id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rs_trip`
--

LOCK TABLES `rs_trip` WRITE;
/*!40000 ALTER TABLE `rs_trip` DISABLE KEYS */;
INSERT INTO `rs_trip` VALUES (1,1001,6,106,'2019-08-22 11:48:34','2019-08-22 11:50:36',12694848,98,36.069094,120.430808,36.06978,120.430129,4,'2019-08-22 09:06:33','2019-08-22 11:50:36'),(9,0,12,2001,'2019-08-23 14:34:05','2019-08-23 14:34:17',12694765,0,36.06978,120.430129,36.06978,120.430129,4,'2019-08-22 12:06:37','2019-08-23 14:34:17'),(10,0,15,2001,'2019-08-26 10:22:16','2019-08-26 10:22:44',12694765,0,36.06978,120.430129,36.06978,120.430129,4,'2019-08-26 10:21:30','2019-08-26 10:22:44'),(11,0,16,2001,'2019-08-26 14:49:11','2019-08-26 14:49:25',12694765,0,36.06978,120.430129,36.06978,120.430129,4,'2019-08-26 14:45:05','2019-08-26 14:49:25');
/*!40000 ALTER TABLE `rs_trip` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-08-26 14:59:02
