/*
TODO my.cnf innodb_autoinc_lock_mode=0 필수
*/


create table osu.BEATMAP_PACK
(
    PACK_ID            int                                   not null
        primary key,
    TYPE               varchar(45)                           not null,
    NAME               varchar(1023)                         null,
    CREATOR            varchar(1023)                         null,
    DATE               varchar(45)                           null,
    DOWNLOAD_URL       varchar(2047)                         null,
    SYSTEM_UPDATE_DATE timestamp default current_timestamp() not null on update current_timestamp()
)
    collate = utf8mb3_unicode_ci;

create table osu.BEATMAP_PACK_SETS
(
    PACK_ID            int                                   not null,
    BEATMAPSET_ID      int(1)                                not null,
    SYSTEM_UPDATE_DATE timestamp default current_timestamp() not null on update current_timestamp(),
    primary key (PACK_ID, BEATMAPSET_ID)
)
    collate = utf8mb3_unicode_ci;

create table osu.SEARCH_CACHE_ARTIST
(
    INDEX_KEY     int unsigned      not null,
    BEATMAPSET_ID int               not null,
    TMP           tinyint default 1 not null,
    primary key (INDEX_KEY, BEATMAPSET_ID)
)
    collate = utf8mb3_unicode_ci;

create index SEARCH_CACHE_ARTIST_BEATMAPSET_ID
    on osu.SEARCH_CACHE_ARTIST (BEATMAPSET_ID);

create table osu.SEARCH_CACHE_CREATOR
(
    INDEX_KEY     int unsigned      not null,
    BEATMAPSET_ID int               not null,
    TMP           tinyint default 1 not null,
    primary key (INDEX_KEY, BEATMAPSET_ID)
)
    collate = utf8mb3_unicode_ci;

create index SEARCH_CACHE_CREATOR_BEATMAPSET_ID
    on osu.SEARCH_CACHE_CREATOR (BEATMAPSET_ID);

create table osu.SEARCH_CACHE_OTHER
(
    INDEX_KEY     int unsigned      not null,
    BEATMAPSET_ID int               not null,
    TMP           tinyint default 1 not null,
    primary key (INDEX_KEY, BEATMAPSET_ID)
)
    collate = utf8mb3_unicode_ci;

create index SEARCH_CACHE_OTHER_BEATMAPSET_ID
    on osu.SEARCH_CACHE_OTHER (BEATMAPSET_ID);

create table osu.SEARCH_CACHE_STRING_INDEX
(
    STRING varchar(256) charset utf8mb4 not null,
    ID     int unsigned auto_increment
        primary key,
    TMP    tinyint default 1            not null,
    constraint STRING_UNIQUE
        unique (STRING)
)
    collate = utf8mb3_unicode_ci
    auto_increment = 1346760;

create table osu.SEARCH_CACHE_TAG
(
    INDEX_KEY     int unsigned      not null,
    BEATMAPSET_ID int               not null,
    TMP           tinyint default 1 not null,
    primary key (INDEX_KEY, BEATMAPSET_ID)
)
    collate = utf8mb3_unicode_ci;

create index SEARCH_CACHE_TAG_BEATMAPSET_ID
    on osu.SEARCH_CACHE_TAG (BEATMAPSET_ID);

create table osu.SEARCH_CACHE_TITLE
(
    INDEX_KEY     int unsigned      not null,
    BEATMAPSET_ID int               not null,
    TMP           tinyint default 1 not null,
    primary key (INDEX_KEY, BEATMAPSET_ID)
)
    collate = utf8mb3_unicode_ci;

create index SEARCH_CACHE_TITLE_BEATMAPSET_ID
    on osu.SEARCH_CACHE_TITLE (BEATMAPSET_ID);

create table osu.SERVER_CACHE_K_V_JSON
(
    `KEY` varchar(254)                 not null
        primary key,
    VALUE longtext collate utf8mb4_bin null,
    constraint SERVER_CACHE_K_V_JSON_KEY_uindex
        unique (`KEY`),
    constraint VALUE
        check (json_valid(`VALUE`))
);

create table osu.beatmap
(
    beatmap_id              int(10)                               not null
        primary key,
    beatmapset_id           int(10)                               not null,
    mode                    varchar(6)                            null,
    mode_int                tinyint(1)                            null,
    status                  varchar(9)                            null,
    ranked                  tinyint(1)                            null,
    total_length            int(10)                               null,
    max_combo               int(10)                               null,
    difficulty_rating       decimal(63, 2)                        null,
    version                 varchar(254)                          null,
    accuracy                decimal(63, 2)                        null,
    ar                      decimal(63, 2)                        null,
    cs                      decimal(63, 2)                        null,
    drain                   decimal(63, 2)                        null,
    bpm                     decimal(63, 2)                        null,
    `convert`               tinyint(1)                            null,
    count_circles           int(10)                               null,
    count_sliders           int(10)                               null,
    count_spinners          int(10)                               null,
    deleted_at              datetime                              null,
    hit_length              int(10)                               null,
    is_scoreable            tinyint(1)                            null,
    last_updated            datetime                              null,
    passcount               int(10)                               null,
    playcount               int(10)                               null,
    checksum                varchar(32)                           null,
    user_id                 int(10)                               null,
    SYSTEM_UPDATE_TIMESTAMP timestamp default current_timestamp() null on update current_timestamp()
)
    collate = utf8mb3_unicode_ci;

create index search0
    on osu.beatmap (beatmapset_id);

create index search1
    on osu.beatmap (mode_int);

create index search2
    on osu.beatmap (total_length);

create index search3
    on osu.beatmap (max_combo);

create index search4
    on osu.beatmap (difficulty_rating);

create index search5
    on osu.beatmap (accuracy);

create index search6
    on osu.beatmap (ar);

create index search7
    on osu.beatmap (cs);

create index search8
    on osu.beatmap (drain);

create index search9
    on osu.beatmap (bpm);

create table osu.beatmapset
(
    beatmapset_id                  int(1)                                not null
        primary key,
    artist                         varchar(254)                          null,
    artist_unicode                 varchar(254) charset utf8mb4          null,
    creator                        varchar(254) charset utf8mb4          null,
    favourite_count                int(1)                                null,
    hype_current                   int(1)                                null,
    hype_required                  int(1)                                null,
    nsfw                           tinyint(1)                            null,
    play_count                     int(1)                                null,
    source                         varchar(254) charset utf8mb4          null,
    status                         varchar(9) charset utf8mb4            null,
    title                          varchar(254) charset utf8mb4          null,
    title_unicode                  varchar(254) charset utf8mb4          null,
    user_id                        int(1)                                null,
    video                          tinyint(1)                            null,
    availability_download_disabled tinyint(1)                            null,
    availability_more_information  text charset utf8mb4                  null,
    bpm                            decimal(63, 2)                        null,
    can_be_hyped                   tinyint(1)                            null,
    discussion_enabled             tinyint(1)                            null,
    discussion_locked              tinyint(1)                            null,
    is_scoreable                   tinyint(1)                            null,
    last_updated                   datetime                              null,
    legacy_thread_url              varchar(254) charset utf8mb4          null,
    nominations_summary_current    int(1)                                null,
    nominations_summary_required   int(1)                                null,
    ranked                         tinyint(1)                            null,
    ranked_date                    datetime                              null,
    storyboard                     tinyint(1)                            null,
    submitted_date                 datetime                              null,
    tags                           text charset utf8mb4                  null,
    has_favourited                 tinyint(1)                            null,
    description                    text charset utf8mb4                  null,
    genre_id                       int(1)                                null,
    genre_name                     varchar(254) charset utf8mb4          null,
    language_id                    int(1)                                null,
    language_name                  varchar(254) charset utf8mb4          null,
    ratings                        varchar(254)                          null,
    SYSTEM_UPDATE_TIMESTAMP        timestamp default current_timestamp() null on update current_timestamp()
)
    collate = utf8mb3_unicode_ci;

create index search106
    on osu.beatmapset (ranked_date);

create index search107
    on osu.beatmapset (favourite_count);

create index search108
    on osu.beatmapset (play_count);

create index search109
    on osu.beatmapset (last_updated);

create index search110
    on osu.beatmapset (title);

create index search112
    on osu.beatmapset (artist);

create index search113
    on osu.beatmapset (ranked_date);

create index search4
    on osu.beatmapset (beatmapset_id);

create index search5
    on osu.beatmapset (ranked);

create index search6
    on osu.beatmapset (nsfw);

create index search7
    on osu.beatmapset (video);

create index search8
    on osu.beatmapset (storyboard);

