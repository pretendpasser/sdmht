CREATE TABLE IF NOT EXISTS `god` (
  `id` bigint NOT NULL,
  `name` varchar(50) NOT NULL,

  `rarity` int(16) NOT NULL,
  `affiliate` int(16) NOT NULL,

  `health` int(16) unsigned NOT NULL,
  `attack` int(16) unsigned NOT NULL,
  `move`  int(16) unsigned NOT NULL,
  `defend` int(16) unsigned NOT NULL,
  `skill_name` varchar(128) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='table of god';


INSERT INTO `god` (`id`, `name`, `rarity`,`affiliate`, `health`,`attack`,`move`,`defend`, `skill_name`)VALUES 
(10,  "奥丁"            ,0,1,   50,3,2,15,    ""),
(11,  "奥丁"            ,0,1,   50,3,2,15,    "德罗普尼尔;机略"),
(12,  "奥丁"            ,0,1,   50,3,2,15,    "冈格尼尔;爆风"),
(13,  "奥丁"            ,0,1,   50,3,2,15,    "英灵指挥;强横"),
(20,  "天照"            ,0,2,   50,3,2,15,    ""),
(21,  "天照"            ,0,2,   50,3,2,15,    "天之岩户;抗拒"),
(22,  "天照"            ,0,2,   50,3,2,15,    "日耀;天明"),
(23,  "天照"            ,0,2,   50,3,2,15,    "八咫镜;反射"),
(30,  "伊米尔"          ,0,3,   50,3,2,15,    ""),
(31,  "伊米尔"          ,0,3,   50,3,2,15,    "战争咆哮;狂风"),
(32,  "伊米尔"          ,0,3,   50,3,2,15,    "暴风雪援军;尖兵"),
(33,  "伊米尔"          ,0,3,   50,3,2,15,    "蛮荒的血脉;血迹"),
(40,  "道德天尊"        ,0,4,   50,3,2,15,    ""),
(41,  "道德天尊"        ,0,4,   50,3,2,15,    "天变地异;天道"),
(42,  "道德天尊"        ,0,4,   50,3,2,15,    "八卦炉;炼化"),
(43,  "道德天尊"        ,0,4,   50,3,2,15,    "变易;流转"),
(50,  "毗湿奴"          ,0,5,   50,3,2,15,    ""),
(51,  "毗湿奴"          ,0,5,   50,3,2,15,    "神龟俱利摩;明灯"),
(52,  "毗湿奴"          ,0,5,   50,3,2,15,    "人狮那罗辛哈;狂乱"),
(53,  "毗湿奴"          ,0,5,   50,3,2,15,    "白马迦尔吉;吉兆"),
(60,  "拉"              ,0,6,   50,3,2,15,    ""),
(61,  "拉"              ,0,6,   50,3,2,15,    "复仇之眼;复仇"),
(62,  "拉"              ,0,6,   50,3,2,15,    "太阳船;秘术"),
(63,  "拉"              ,0,6,   50,3,2,15,    "公理天秤;太阳神"),
(70,  "宙斯"            ,0,7,   50,3,2,15,    ""),
(71,  "宙斯"            ,0,7,   50,3,2,15,    "无上权能;君王"),
(72,  "宙斯"            ,0,7,   50,3,2,15,    "威压雷霆;君王"),
(73,  "宙斯"            ,0,7,   50,3,2,15,    "弑父霸王;君王"),
(80,  "元始天尊"        ,0,8,   50,3,2,15,    ""),
(81,  "元始天尊"        ,0,8,   50,3,2,15,    "天;天"),
(82,  "元始天尊"        ,0,8,   50,3,2,15,    "乾;乾"),
(83,  "元始天尊"        ,0,8,   50,3,2,15,    "苍;苍"),

(1000, "庄周凰之型"      ,2,0,   50,3,1,15,     "庄周凰之型"),
(1001, "庄周鹏之型"      ,3,0,   36,1,0,15,     "风鹏"),
(1002, "庄周鲲之型"      ,3,0,   40,2,1,15,     "深渊"),
(1003, "南华整备士"      ,1,0,   40,3,1,15,     "整备"),
(1004, "夜之蝶"          ,1,0,    1,0,1,15,     "夜之蝶"),
(1005, "南华管理者"      ,1,0,   40,4,1,15,     "偶像制造"),

(101, "史基尼尔"        ,1,1,   35,7,1,15,     "信使"),
(102, "弗丽嘉"          ,1,1,   32,6,1,15,     "爱之援助"),
(103, "弗雷"            ,1,1,   24,1,1,15,     "丰饶之路"),
(104, "霍德尔"          ,2,1,   35,5,1,15,     "暗闇"),
(105, "西芙"            ,2,1,   35,4,1,15,     "大地护佑"),
(106, "海拉"            ,2,1,   35,9,1,15,     "冥河幽深"),
(107, "海姆达尔"        ,2,1,   40,2,1,15,     "彩虹桥守卫"),
(108, "洛基"            ,2,1,   32,3,1,15,     "欺诈者"),
(109, "提尔"            ,2,1,   32,5,1,15,     "英灵引导"),
(110, "布拉吉"          ,2,1,   32,5,1,15,     "诗与乐"),
(111, "古尔薇格"        ,2,1,   40,3,1,15,     "煽动者"),
(112, "博德"            ,3,1,   32,5,1,15,     "光辉之子"),
(113, "尼约德"          ,3,1,   36,2,1,15,     "盛夏"),
(114, "希格德莉法"      ,3,1,   35,8,1,15,     "胜利宣判者"),
(115, "托尔"            ,3,1,   40,6,1,15,     "雷鸣震慑"),
(116, "斯卡蒂"          ,3,1,   33,5,1,15,     "凛冬"),
(117, "芙蕾雅"          ,3,1,   40,4,1,15,     "永驻的青春"),
(118, "霍德尔-常暗"     ,3,1,   40,0,1,15,     "常暗"),

(201, "大和武尊"        ,1,2,   35,5,1,15,     "奉天承命"),
(202, "天之常立"        ,1,2,   32,4,1,15,     "常在战阵"),
(203, "建御雷"          ,1,2,   36,7,1,15,     "先锋"),
(204, "思兼神"          ,1,2,   35,2,1,15,     "智囊"),
(205, "事代主"          ,2,2,   45,0,1,15,     "调解者"),
(206, "大国主"          ,2,2,   25,6,1,15,     "国仇"),
(207, "天手力男"        ,2,2,   50,0,1,15,     "豪腕"),
(208, "天探女"          ,2,2,   50,5,1,15,     "陪臣"),
(209, "天若日子"        ,2,2,   50,5,1,15,     "良禽"),
(210, "建御名方"        ,2,2,   35,5,1,15,     "赤口之祟"),
(211, "水蛭子"          ,2,2,   35,4,1,15,     "暴发户"),
(212, "火之迦具土"      ,2,2,   32,3,1,15,     "业火"),
(213, "经津主"          ,2,2,   32,5,1,15,     "剑锋磨砺"),
(214, "荒霸吐"          ,2,2,   35,4,1,15,     "昔日的荣光"),
(215, "冥后奇稻田"      ,3,2,   28,3,1,15,     "丰壤的冥后"),
(216, "大物主"          ,3,2,   40,4,1,15,     "国津魂"),
(217, "少彦名"          ,3,2,   35,3,1,15,     "往来常世"),
(218, "月读"            ,3,2,   40,3,1,15,     "悖逆"),
(219, "月读十六夜"      ,3,2,   39,5,1,15,     "真月"),
(220, "毗沙门天"        ,3,2,   36,8,1,15,     "军神"),
(221, "辉夜"            ,3,2,   40,4,1,15,     "燕之子安贝"),
(222, "须佐之男"        ,3,2,   35,5,1,15,     "不羁"),
(223, "须佐之男•烬"     ,3,2,   45,0,1,15,     "天丛云"),

(301, "哈提"            ,1,3,   40,6,1,15,     "吞食"),
(302, "奥尔布达"        ,1,3,   30,3,1,15,     "魔兽之母"),
(303, "梅妮亚"          ,1,3,   60,4,1,15,     "贪婪之祸"),
(304, "菲妮亚"          ,1,3,   60,4,1,15,     "强欲之灾"),
(305, "法夫尼尔"        ,1,3,   36,3,1,15,     "龙之宝藏"),
(306, "索列姆"          ,1,3,   35,8,1,15,     "荒土的蛮王"),
(307, "蒂阿兹"          ,1,3,   40,4,1,15,     "隼魔"),
(308, "赫拉斯瓦尔"      ,1,3,   40,8,1,15,     "追魂魔喙"),
(309, "斯库鲁"          ,2,3,   38,5,1,15,     "吞噬太阳"),
(310, "约莫加德"        ,2,3,   35,5,1,15,     "缠绕世界"),
(311, "芬里尔"          ,2,3,   40,5,1,15,     "吞噬世界"),
(312, "赫朗格尼尔"      ,2,3,   50,0,1,15,     "勇壮豪伟"),
(313, "【洛基】"        ,3,3,   40,6,1,15,     "血之亲族"),
(314, "乌特加洛基"      ,3,3,   40,5,1,15,     "魔性"),
(315, "史尔特"          ,3,3,   25,7,1,15,     "红王"),
(316, "冰之洛基"        ,3,3,   32,4,1,15,     "寒冰女王"),
(317, "尼德霍格"        ,3,3,   60,2,1,15,     "至黑魔龙"),
(318, "摩克卡菲"        ,3,3,   60,12,1,15,    "穷鼠之噬"),
(319, "蛇魔约莫加德"    ,3,3,   38,5,1,15,     "世界之毒"),
(320, "阿尔贝莉琪"      ,3,3,   20,4,1,15,     "咒缚之指环"),

(401, "嫦娥重装型"      ,3,4,   36,4,1,15,     "要塞广寒"),
(402, "嫦娥"            ,1,4,   40,5,1,15,    "守月卫士"),
(403, "太阴星君"        ,1,4,   32,4,1,15,    "月崩"),
(404, "太白金星"        ,1,4,   40,4,1,15,    "启明"),
(405, "太乙真人"        ,1,4,   45,2,1,15,    "金眼"),
(406, "通玄真人"        ,1,4,   32,4,1,15,    "微言"),
(407, "哪吒"            ,2,4,   35,6,1,15,    "万宝"),
(408, "太公望"          ,2,4,   34,6,1,15,    "八卦图"),
(409, "李靖"            ,2,4,   40,4,1,15,    "托塔天王"),
(410, "杨戬"            ,2,4,   8,1,1,15,     "天眼招来"),
(411, "冲虚真人"        ,1,4,   34,4,1,15,    "贵虚"),
(412, "王诩"            ,3,4,   35,4,1,15,    "鬼谷神算"),
(413, "哪吒原初型"      ,3,4,   35,2,1,15,    "轻装上阵"),
(414, "孙悟空"          ,3,4,   35,6,1,15,    "大闹天宫"),
(415, "日耀帝君"        ,3,4,   38,6,1,15,    "天狱"),
(416, "洞灵真人"        ,3,4,   32,8,1,15,    "变化万端"),
(417, "后羿"            ,3,4,   35,4,1,15,    "狩日射手"),
(418, "仲尼"            ,3,4,   45,8,1,15,    "天之木铎"),

(501, "拉克西米"       ,1,5,   32,4,1,15,     "莲花"),
(502, "娜迦"           ,1,5,   35,5,1,15,     "紫渊之民"),
(503, "因陀罗"         ,2,5,   40,5,1,15,     "威严万丈"),
(504, "增长天"         ,2,5,   35,3,1,15,     "万物赐福"),
(505, "多闻天"         ,2,5,   35,5,1,15,     "军神"),
(506, "广目天"         ,2,5,   35,5,1,15,     "远见"),
(507, "持国天"         ,2,5,   35,5,1,15,     "天门守护"),
(508, "犍尼萨"         ,2,5,   45,0,1,15,     "象的智略"),
(509, "哪吒俱伐罗"     ,3,5,   37,3,1,15,     "除魔卫道"),
(510, "哈努曼"         ,3,5,   40,9,1,15,     "猿军"),
(511, "帕尔瓦蒂"       ,3,5,   40,5,1,15,     "迦梨"),
(512, "湿婆"           ,3,5,   25,5,1,15,     "诸界破灭"),
(513, "罗睺"           ,3,5,   25,7,1,15,     "凶星"),
(514, "迦楼罗"         ,3,5,   35,6,1,15,     "苍空之民"),
(515, "鸠摩罗"         ,3,5,   48,5,1,15,     "战神之尊"),
(516, "吉祥天女"       ,4,5,   54,0,1,15,     "第三宝"),

(601, "凯布利"          ,1,6,   22,8,1,15,     "日轮疾驰"),
(602, "阿穆特"          ,2,6,   45,3,1,15,     "啃噬"),
(603, "阿努比斯"        ,2,6,   35,6,1,15,     "死灵守护者"),
(604, "索贝克"          ,2,6,   45,0,1,15,     "龙颚的巨兽"),
(605, "托特"            ,2,6,   32,5,1,15,     "智慧回响"),
(606, "姆特"            ,2,6,   45,5,1,15,     "战争呼唤"),
(607, "奥西里斯"        ,2,6,   28,8,1,15,     "天命循环"),
(608, "塞赫迈特"        ,2,6,   45,6,1,15,     "母狮"),
(609, "伊西斯"          ,2,6,   40,3,1,15,     "生命永恒"),
(610, "亚顿"            ,2,6,   35,4,1,15,     "古之朝阳"),
(611, "塞特"            ,3,6,   40,8,1,15,     "祸土"),
(612, "奈芙蒂斯"        ,3,6,   36,2,1,15,     "冥土之盾"),
(613, "孔斯"            ,3,6,   40,2,1,15,     "月之诅咒"),
(614, "弗雷·夏之王"     ,3,6,   28,3,1,15,     "华纳神王"),
(615, "荷鲁斯"          ,3,6,   43,4,1,15,     "复仇之眼"),
(616, "贝斯特"          ,3,6,   40,4,1,15,     "吉兆之猫"),

(701, "赫菲斯托斯"      ,1,7,   33,0,1,15,     "神匠"),
(702, "阿特拉斯"        ,1,7,   30,5,1,15,     "暴政囚徒"),
(703, "赫拉"            ,2,7,   40,4,1,15,     "爱之形"),
(704, "阿尔忒弥斯"      ,2,7,   45,4,1,15,     "疯乱"),
(705, "阿波罗"          ,2,7,   40,3,1,15,     "偶像"),
(706, "阿瑞斯"          ,2,7,   40,3,1,15,     "战火延烧"),
(707, "哈迪斯"          ,3,7,   30,4,1,15,     "不可逆之路"),
(708, "堤丰"            ,3,7,   42,5,1,15,     "污秽血脉"),
(709, "波塞冬"          ,3,7,   38,5,1,15,     "绝海神殿"),
(710, "涅墨西斯"        ,3,7,   40,5,1,15,     "复仇的才能"),
(711, "雅典娜"          ,3,7,   35,2,1,15,     "奇谋"),
(712, "普罗米修斯"      ,4,7,   35,4,1,15,     "盗火者"),
(713, "埃庇米修斯"      ,4,7,   45,3,1,15,     "灾难之盒"),
(714, "狄俄尼索斯"      ,4,7,   55,0,1,15,     "狂欢的盛宴"),

(801, "共工"            ,2,8,   35,4,1,15,     "天理崩塌"),
(802, "列御寇"          ,2,8,   38,2,1,15,     "大道冲虚"),
(803, "太乙救苦天尊"    ,3,8,   40,4,1,15,     "救苦度亡"),
(804, "应龙"            ,3,8,   40,7,1,15,     "水"),
(805, "庚桑楚"          ,3,8,   38,8,1,15,     "合真"),
(806, "旱魃"            ,3,8,   44,5,1,15,     "节制"),
(807, "浑沌"            ,3,8,   10,5,1,15,     "窍"),
(808, "烛阴"            ,3,8,   45,5,1,15,     "清浊之辩"),
(809, "相柳"            ,3,8,   35,4,1,15,     "毒潭"),
(810, "织女"            ,3,8,   44,3,1,15,     "云中星"),
(811, "辛计然"          ,3,8,   40,3,1,15,     "大道通玄")
;