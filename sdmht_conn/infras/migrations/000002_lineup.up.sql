CREATE TABLE IF NOT EXISTS `lineup` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `account_id` bigint unsigned NOT NULL,
  `name` varchar(128) NOT NULL,
  `units` varchar(256) NOT NULL,
  `card_library` varchar(256) NOT NULL DEFAULT "",
  `enabled` tinyint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='table of lineup';
