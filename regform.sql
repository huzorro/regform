-- MySQL dump 10.13  Distrib 5.5.41, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: regform
-- ------------------------------------------------------
-- Server version	5.5.41-0ubuntu0.14.04.1

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
-- Table structure for table `regform_log`
--

DROP TABLE IF EXISTS `regform_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `regform_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `info` varchar(1000) NOT NULL DEFAULT '',
  `type` int(11) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  `logtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `regform_log`
--

LOCK TABLES `regform_log` WRITE;
/*!40000 ALTER TABLE `regform_log` DISABLE KEYS */;
INSERT INTO `regform_log` VALUES (11,'{\"Id\":11,\"University\":\"发达东方123456\",\"Level\":\"adf\",\"Leader\":\"adf\",\"LeaderGender\":\"asdf\",\"LeaderContact\":\"adsf\",\"Pm\":\"adsf\",\"PmGender\":\"asdf\",\"PmContact\":\"asdf\",\"Topic\":\"asdf\",\"Intro\":\"撒地方啊但是发范德萨123\",\"Logtime\":\"20150424\"}',1,2,'2015-04-24 07:40:21'),(15,'{\"Id\":15,\"Name\":\"大众\",\"Gender\":\"男\",\"Birth\":\"asdf\",\"Believe\":\"asdf\",\"Unit\":\"asdf\",\"Mobile\":\"adf\",\"Email\":\"asdf\",\"Positions\":\"adf\",\"Director\":\"asdf\",\"Evaluation\":\"啊短发的辐射度方法大啊神烦大叔东方闪电\",\"Logtime\":\"20150424\"}',2,2,'2015-04-24 07:57:10');
/*!40000 ALTER TABLE `regform_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sp_access_privilege`
--

DROP TABLE IF EXISTS `sp_access_privilege`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sp_access_privilege` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pri_group` varchar(500) NOT NULL DEFAULT '' COMMENT '1;2;3;4;5',
  `pri_rule` int(11) NOT NULL DEFAULT '0' COMMENT '1:all, 2:allow, 4:ban',
  `logtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sp_access_privilege`
--

LOCK TABLES `sp_access_privilege` WRITE;
/*!40000 ALTER TABLE `sp_access_privilege` DISABLE KEYS */;
INSERT INTO `sp_access_privilege` VALUES (1,'',1,'2015-04-24 07:02:08'),(2,'1',2,'2015-04-24 07:02:08'),(3,'1',4,'2015-04-24 07:02:09');
/*!40000 ALTER TABLE `sp_access_privilege` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sp_menu_template`
--

DROP TABLE IF EXISTS `sp_menu_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sp_menu_template` (
  `id` int(11) NOT NULL DEFAULT '0' COMMENT '1 2 4 8',
  `title` varchar(100) NOT NULL DEFAULT '' COMMENT '团体报名表 个人报名表',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT 'team, person',
  `logtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sp_menu_template`
--

LOCK TABLES `sp_menu_template` WRITE;
/*!40000 ALTER TABLE `sp_menu_template` DISABLE KEYS */;
INSERT INTO `sp_menu_template` VALUES (1,'团体报名表','team','2015-04-24 07:02:24'),(2,'个人报名表','person','2015-04-24 07:02:24');
/*!40000 ALTER TABLE `sp_menu_template` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sp_node_privilege`
--

DROP TABLE IF EXISTS `sp_node_privilege`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sp_node_privilege` (
  `id` int(11) NOT NULL DEFAULT '0',
  `name` varchar(100) NOT NULL DEFAULT '',
  `node` varchar(500) NOT NULL DEFAULT '' COMMENT '1:/login, 2:/rLogin, 4:/logout, 8:/regform/team, 16:/regform/person, 32:/regform/team/submit, 64:/regform/person/submit, 128:/regform/admin/team, 256:/regform/admin/person, 512:/getform, 1024:/updateform, 2048:/deleteform',
  `logtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sp_node_privilege`
--

LOCK TABLES `sp_node_privilege` WRITE;
/*!40000 ALTER TABLE `sp_node_privilege` DISABLE KEYS */;
INSERT INTO `sp_node_privilege` VALUES (1,'登录页','/login','2015-04-24 07:01:42'),(2,'登录验证请求','/rLogin','2015-04-24 07:01:42'),(4,'退出登录','/logout','2015-04-24 14:03:51'),(8,'团队报名表单','/regform/team','2015-04-24 07:01:42'),(16,'个人报名表单','/regform/person','2015-04-24 07:01:42'),(32,'提交团队报名表单','/regform/team/submit','2015-04-24 07:01:42'),(64,'提交个人报名表单','/regform/person/submit','2015-04-24 07:01:42'),(128,'团队报名表单管理','/regform/admin/team','2015-04-24 07:01:42'),(256,'个人报名表单管理','/regform/admin/person','2015-04-24 07:01:42'),(512,'获取单个表单','/getform','2015-04-24 07:01:42'),(1024,'更新单个表单','/updateform','2015-04-24 07:01:42'),(2048,'删除单个表单','/deleteform','2015-04-24 07:01:42');
/*!40000 ALTER TABLE `sp_node_privilege` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sp_role`
--

DROP TABLE IF EXISTS `sp_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sp_role` (
  `id` int(11) NOT NULL DEFAULT '0',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT 'user, services, admin, guess',
  `privilege` int(11) NOT NULL DEFAULT '0',
  `menu` int(11) NOT NULL DEFAULT '0',
  `logtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sp_role`
--

LOCK TABLES `sp_role` WRITE;
/*!40000 ALTER TABLE `sp_role` DISABLE KEYS */;
INSERT INTO `sp_role` VALUES (1,'匿名用户',123,0,'2015-04-24 07:21:49'),(2,'管理员',4095,3,'2015-04-24 07:01:25');
/*!40000 ALTER TABLE `sp_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sp_user`
--

DROP TABLE IF EXISTS `sp_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sp_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(100) NOT NULL DEFAULT '',
  `password` varchar(100) NOT NULL DEFAULT '',
  `roleid` int(11) NOT NULL DEFAULT '0',
  `accessid` int(11) NOT NULL DEFAULT '0',
  `logtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sp_user`
--

LOCK TABLES `sp_user` WRITE;
/*!40000 ALTER TABLE `sp_user` DISABLE KEYS */;
INSERT INTO `sp_user` VALUES (1,'root','admin',2,1,'2015-04-24 07:00:42');
/*!40000 ALTER TABLE `sp_user` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2015-04-24 22:24:40
