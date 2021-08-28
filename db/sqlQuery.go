package db


const QueryBeatmap = `select * from osu.beatmap where beatmapset_id in( %s ) order by difficulty_rating desc;`

const QuerySearchBeatmapSet = `
select * from osu.beatmapset 
	where ranked in( %s ) 
	AND
	beatmapset_id in (select distinct beatmapset_id from osu.beatmap where ranked in( %s ) AND mode_int in ( %s ) ) 
order by %s %s ;
`

const QuerySearchBeatmapSetWhitQueryText = `
select * from osu.beatmapset 
where  ranked in( %s ) 
AND beatmapset_id in (select distinct beatmapset_id from osu.beatmap where ranked in( %s ) AND mode_int in ( %s ) ) 
AND beatmapset_id in (select beatmapset_id from osu.search_index where MATCH(text) AGAINST(?))
order by %s %s ;
`

//const QueryAPILog = `INSERT INTO osu.api_log (time, request_id, remote_ip, host, method, uri, user_agent, status, error, latency, latency_human, bytes_in, bytes_out) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?);`

const UpsertBeatmap = `
INSERT INTO osu.beatmap
(
	beatmap_id,beatmapset_id,mode,mode_int,status,	ranked,total_length,max_combo,difficulty_rating,version,
	accuracy,ar,cs,drain,bpm,`+"`convert`"+`,count_circles,count_sliders,count_spinners,deleted_at,
	hit_length,is_scoreable,last_updated,passcount,playcount,	checksum,user_id
)VALUES(
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?
)ON DUPLICATE KEY UPDATE 
	beatmapset_id = VALUES(beatmapset_id), mode = VALUES(mode), mode_int = VALUES(mode_int), status = VALUES(status), 
	ranked = VALUES(ranked), total_length = VALUES(total_length), max_combo = VALUES(max_combo), difficulty_rating = VALUES(difficulty_rating), 
	version = VALUES(version), 	accuracy = VALUES(accuracy), ar = VALUES(ar), cs = VALUES(cs), drain = VALUES(drain), bpm = VALUES(bpm), 
	`+"`convert` = VALUES(`convert`"+`), count_circles = VALUES(count_circles), count_sliders = VALUES(count_sliders), 
	count_spinners = VALUES(count_spinners), deleted_at = VALUES(deleted_at), 	hit_length = VALUES(hit_length), 
	is_scoreable = VALUES(is_scoreable), last_updated = VALUES(last_updated), passcount = VALUES(passcount), playcount = VALUES(playcount), 
	checksum = VALUES(checksum), user_id = VALUES(user_id)
;
`

const UpsertBeatmapSet = `
INSERT INTO osu.beatmapset(
	beatmapset_id,artist,artist_unicode,creator,favourite_count,
	hype_current,hype_required,nsfw,play_count,source,
	status,title,title_unicode,user_id,video,
	availability_download_disabled,availability_more_information,bpm,can_be_hyped,discussion_enabled,
	discussion_locked,is_scoreable,last_updated,legacy_thread_url,nominations_summary_current,
	nominations_summary_required,ranked,ranked_date,storyboard,submitted_date,
	tags,has_favourited,description,genre_id,genre_name,
	language_id,language_name,ratings
)VALUES(
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?,?
)ON DUPLICATE KEY UPDATE 
	artist= VALUES(artist), artist_unicode= VALUES(artist_unicode), creator= VALUES(creator), favourite_count= VALUES(favourite_count), 
	hype_current= VALUES(hype_current), hype_required= VALUES(hype_required), nsfw= VALUES(nsfw), play_count= VALUES(play_count), source= VALUES(source), 
	status= VALUES(status), title= VALUES(title), title_unicode= VALUES(title_unicode), user_id= VALUES(user_id), video= VALUES(video), 
	availability_download_disabled= VALUES(availability_download_disabled), availability_more_information= VALUES(availability_more_information), 
	bpm= VALUES(bpm), can_be_hyped= VALUES(can_be_hyped), discussion_enabled= VALUES(discussion_enabled), 
	discussion_locked= VALUES(discussion_locked), is_scoreable= VALUES(is_scoreable), last_updated= VALUES(last_updated), legacy_thread_url= VALUES(legacy_thread_url), 
	nominations_summary_current= VALUES(nominations_summary_current), 	nominations_summary_required= VALUES(nominations_summary_required), 
	ranked= VALUES(ranked), ranked_date= VALUES(ranked_date), storyboard= VALUES(storyboard), submitted_date= VALUES(submitted_date), 
	tags= VALUES(tags), has_favourited= VALUES(has_favourited), description= VALUES(description), genre_id= VALUES(genre_id), genre_name= VALUES(genre_name), 
	language_id= VALUES(language_id), language_name= VALUES(language_name), ratings= VALUES(ratings)
;
`


const GetDownloadBeatmapData = `SELECT beatmapset_id,artist,title,last_updated,video FROM osu.beatmapset WHERE beatmapset_id = ?`


/*
CREATE TABLE `beatmap` (
  `beatmap_id` int(10) NOT NULL,
  `beatmapset_id` int(10) NOT NULL,
  `mode` varchar(6) COLLATE utf8_unicode_ci DEFAULT NULL,
  `mode_int` tinyint(1) DEFAULT NULL,
  `status` varchar(9) COLLATE utf8_unicode_ci DEFAULT NULL,
  `ranked` tinyint(1) DEFAULT NULL,
  `total_length` int(10) DEFAULT NULL,
  `max_combo` int(10) DEFAULT NULL,
  `difficulty_rating` decimal(63,2) DEFAULT NULL,
  `version` varchar(254) COLLATE utf8_unicode_ci DEFAULT NULL,
  `accuracy` decimal(63,2) DEFAULT NULL,
  `ar` decimal(63,2) DEFAULT NULL,
  `cs` decimal(63,2) DEFAULT NULL,
  `drain` decimal(63,2) DEFAULT NULL,
  `bpm` decimal(63,2) DEFAULT NULL,
  `convert` tinyint(1) DEFAULT NULL,
  `count_circles` int(10) DEFAULT NULL,
  `count_sliders` int(10) DEFAULT NULL,
  `count_spinners` int(10) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `hit_length` int(10) DEFAULT NULL,
  `is_scoreable` tinyint(1) DEFAULT NULL,
  `last_updated` datetime DEFAULT NULL,
  `passcount` int(10) DEFAULT NULL,
  `playcount` int(10) DEFAULT NULL,
  `checksum` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `user_id` int(10) DEFAULT NULL,
  PRIMARY KEY (`beatmap_id`),
  UNIQUE KEY `beatmap_id_UNIQUE` (`beatmap_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


CREATE TABLE `beatmapset` (
  `beatmapset_id` int(1) NOT NULL,
  `artist` varchar(254) COLLATE utf8_unicode_ci DEFAULT NULL,
  `artist_unicode` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `creator` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `favourite_count` int(1) DEFAULT NULL,
  `hype_current` int(1) DEFAULT NULL,
  `hype_required` int(1) DEFAULT NULL,
  `nsfw` tinyint(1) DEFAULT NULL,
  `play_count` int(1) DEFAULT NULL,
  `source` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `status` varchar(9) CHARACTER SET utf8mb4 DEFAULT NULL,
  `title` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `title_unicode` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `user_id` int(1) DEFAULT NULL,
  `video` tinyint(1) DEFAULT NULL,
  `availability_download_disabled` tinyint(1) DEFAULT NULL,
  `availability_more_information` text CHARACTER SET utf8mb4 DEFAULT NULL,
  `bpm` decimal(63,2) DEFAULT NULL,
  `can_be_hyped` tinyint(1) DEFAULT NULL,
  `discussion_enabled` tinyint(1) DEFAULT NULL,
  `discussion_locked` tinyint(1) DEFAULT NULL,
  `is_scoreable` tinyint(1) DEFAULT NULL,
  `last_updated` datetime DEFAULT NULL,
  `legacy_thread_url` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `nominations_summary_current` int(1) DEFAULT NULL,
  `nominations_summary_required` int(1) DEFAULT NULL,
  `ranked` tinyint(1) DEFAULT NULL,
  `ranked_date` datetime DEFAULT NULL,
  `storyboard` tinyint(1) DEFAULT NULL,
  `submitted_date` datetime DEFAULT NULL,
  `tags` text CHARACTER SET utf8mb4 DEFAULT NULL,
  `has_favourited` tinyint(1) DEFAULT NULL,
  `description` text CHARACTER SET utf8mb4 DEFAULT NULL,
  `genre_id` int(1) DEFAULT NULL,
  `genre_name` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `language_id` int(1) DEFAULT NULL,
  `language_name` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `ratings` varchar(254) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`beatmapset_id`),
  UNIQUE KEY `id_UNIQUE` (`beatmapset_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE `search_index` (
  `beatmapset_id` int(11) NOT NULL,
  `text` varchar(1023) DEFAULT '',
  PRIMARY KEY (`beatmapset_id`),
  FULLTEXT KEY `text` (`text`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 트리거================
DELIMITER $$
CREATE DEFINER=`root`@`%` TRIGGER `osu`.`beatmapset_BEFORE_INSERT` BEFORE INSERT ON `beatmapset` FOR EACH ROW
BEGIN
	INSERT INTO `osu`.`search_index`(`beatmapset_id`,`text`)
    VALUES(
		NEW.beatmapset_id,
		concat_ws(' ',NEW.beatmapset_id,NEW.artist,NEW.creator,NEW.title)
    )ON DUPLICATE KEY UPDATE text = concat_ws(' ',NEW.beatmapset_id,NEW.artist,NEW.creator,NEW.title);
END$$
DELIMITER ;
-- 트리거================
DELIMITER $$
CREATE DEFINER=`root`@`%` TRIGGER `osu`.`beatmapset_BEFORE_UPDATE` BEFORE UPDATE ON `beatmapset` FOR EACH ROW
BEGIN
	INSERT INTO `osu`.`search_index`(`beatmapset_id`,`text`)
    VALUES(
		NEW.beatmapset_id,
		concat_ws(' ',NEW.beatmapset_id,NEW.artist,NEW.creator,NEW.title)
    )ON DUPLICATE KEY UPDATE text = concat_ws(' ',NEW.beatmapset_id,NEW.artist,NEW.creator,NEW.title);
END$$
DELIMITER ;
-- 트리거================

*/