namespace go leaf

enum Error {
    SEG_SERVER_UNSPECIFIED = 0;
    BIZ_TAG_NOT_FOUND = 1;
    DB_ERROR = 2;
    ID_CACHE_INIT_FAILED = 3;
    ID_TWO_SEGMENTS_ARE_NULL = 4;

    // 时钟回拨超过5ms
    SNOWFLAKE_TIME_EXCEPTION = 5;
    SNOWFLAKE_ID_ILLEGAL = 6;
}

// segment

struct IdRequest {
    1: string tag;
}

struct IdResponse {
    1: string id;
}

// snowflake

struct DecodeSnokflakeRequest {
    1: string id;
}

struct DecodeSnokflakeResponse {
    1: string timestamp;
    2: string sequenceId;
    3: string workerId;
}

service LeafService {
    // 号段模式
    IdResponse GenSegmentId(1: IdRequest request) (api.post = "v1/segment");
    // 雪花ID
    IdResponse GenSnowflakeId(1: IdRequest request) (api.post = "v1/snowflake")
    // 解析雪花ID
    DecodeSnokflakeResponse DecodeSnowflakeId(1: DecodeSnokflakeRequest request) (api.post = "v1/snowflake/decode");
}