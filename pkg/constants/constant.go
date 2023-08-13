package constants

// connection information
const (
	// DB
	MySqlDefaultDSN = "douyin:douyin123@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"

	// Minio
	MinioDefaultEndpoint  = "lcoalhost:9000"
	MinioAccessKeyID      = "douyin"
	MinioSecrectAccessKey = "douyin123"
	MiniouseSSL           = false

	// Redis
	RedisDefaultAddr = "localhost:6379"
	RedisDefaultPwd  = "douyin123"
)

// contants in project
const (
	UserTableName      = "users"
	FollowsTableName   = "follows"
	VideosTableName    = "videos"
	MessageTableName   = "messages"
	FavoritesTableName = "likes"
	CommentsTableName  = "comments"

	VideoFeedCount       = 30
	FavoriteActionType   = 1
	UnFavoriteActionType = 2

	MinioVideoBucketName = "videobucket"
	MinioImageBucketName = "imagebucket"

	TestSign       = "testaccount! offer"
	TestAva        = "avatar/test1.jpg"
	TestBackGround = "background/test1.jpg"
)
