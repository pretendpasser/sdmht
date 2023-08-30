CREATE TABLE IF NOT EXISTS `god` (
  `id` bigint unsigned NOT NULL,
  `name` varchar(50) NOT NULL,

  `rarity` int(16) NOT NULL,
  `affiliate` int(16) NOT NULL,

  `health` int(16) unsigned NOT NULL,
  `attack` int(16) unsigned NOT NULL,
  `move`  int(16) unsigned NOT NULL,
  `defend` int(16) unsigned NOT NULL,

  `nomove` tinyint NOT NULL,
  `noattack` tinyint NOT NULL,
  `nocure` tinyint NOT NULL,

  PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='table of god';


INSERT INTO `god` (`id`, `name`, `rarity`,`affiliate`, `health`,`attack`,`move`,`defend`, `nomove`,`noattack`,`nocure`)VALUES 
(10,  "奥丁-1"            ,0,1,   50,3,2,15,    0,0,0),
(11,  "奥丁-2"            ,0,1,   50,3,2,15,    0,0,0),
(12,  "奥丁-3"            ,0,1,   50,3,2,15,    0,0,0),
(20,  "天照"            ,0,2,   50,3,2,15,    0,0,0),
(30,  "伊米尔"          ,0,3,   50,3,2,15,    0,0,0),
(40,  "道德天尊"        ,0,4,   50,3,2,15,    0,0,0),
(50,  "毗湿奴"          ,0,5,   50,3,2,15,    0,0,0),
(60,  "拉"              ,0,6,   50,3,2,15,    0,0,0),
(70,  "宙斯"            ,0,7,   50,3,2,15,    0,0,0),
(80,  "元始天尊"        ,0,8,   50,3,2,15,    0,0,0),

(401, "嫦娥重装型"      ,3,4,   36,4,1,15,    0,0,0),
(402, "嫦娥"            ,1,4,   40,5,1,15,    0,0,0)

;