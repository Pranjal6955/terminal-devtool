use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct ProcessRequest {
    pub input: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub output: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub resolution: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bitrate: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub format: Option<String>,
    #[serde(skip_serializing_if = "is_false")]
    pub dry_run: bool,
}

fn is_false(b: &bool) -> bool {
    !*b
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ProcessResponse {
    pub output: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CompareRequest {
    pub original: String,
    pub processed: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct MediaInfo {
    pub filename: String,
    pub format: String,
    pub duration: String,
    pub resolution: String,
    pub bitrate: String,
    pub size: u64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CompareResult {
    pub original: MediaInfo,
    pub processed: MediaInfo,
    pub size_diff_percent: f64,
}

#[derive(Debug)]
pub struct CompressOptions {
    pub input: String,
    pub output: Option<String>,
    pub bitrate: String,
}
