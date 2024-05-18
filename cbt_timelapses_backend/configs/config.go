package configs

import "time"

const DIRECTORY_FOLDER_SCRIPT = "/home/blunder/bin/create_video.sh"
const SCREENSHOTS_FOLDER = "/mnt/seefetch"
const ORDER_TTL = 7 * 24 * time.Hour
const PORT_SERVER = ":5000"

const STATUS_OK = 200
const STATUS_WAITING = 300
const STATUS_ERROR = 400
