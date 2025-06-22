use anyhow::{Context, Result};
use reqwest::blocking::Client;
use std::time::Duration;

use crate::models::{CompareRequest, CompareResult, MediaInfo, ProcessRequest, ProcessResponse};

const DEFAULT_BACKEND_URL: &str = "http://localhost:8080";
const TIMEOUT_SECONDS: u64 = 30;

pub struct ApiClient {
    client: Client,
    base_url: String,
}

impl ApiClient {
    pub fn new(base_url: Option<String>) -> Self {
        let url = base_url.unwrap_or_else(|| DEFAULT_BACKEND_URL.to_string());
        
        let client = Client::builder()
            .timeout(Duration::from_secs(TIMEOUT_SECONDS))
            .build()
            .expect("Failed to create HTTP client");
            
        Self {
            client,
            base_url: url,
        }
    }
    
    pub fn process_media(&self, request: ProcessRequest) -> Result<ProcessResponse> {
        let url = format!("{}/api/process", self.base_url);
        
        let response = self.client
            .post(&url)
            .json(&request)
            .send()
            .context("Failed to send process request to backend")?;
            
        if !response.status().is_success() {
            return Err(anyhow::anyhow!(
                "Backend returned error: {} - {}",
                response.status(),
                response.text().unwrap_or_default()
            ));
        }
        
        let result = response
            .json::<ProcessResponse>()
            .context("Failed to parse process response")?;
            
        Ok(result)
    }
    
    pub fn compare_media(&self, original: &str, processed: &str) -> Result<CompareResult> {
        let url = format!("{}/api/compare", self.base_url);
        
        let request = CompareRequest {
            original: original.to_string(),
            processed: processed.to_string(),
        };
        
        let response = self.client
            .post(&url)
            .json(&request)
            .send()
            .context("Failed to send compare request to backend")?;
            
        if !response.status().is_success() {
            return Err(anyhow::anyhow!(
                "Backend returned error: {} - {}",
                response.status(),
                response.text().unwrap_or_default()
            ));
        }
        
        let result = response
            .json::<CompareResult>()
            .context("Failed to parse compare response")?;
            
        Ok(result)
    }
    
    pub fn get_media_info(&self, file_path: &str) -> Result<MediaInfo> {
        let url = format!("{}/api/info?path={}", self.base_url, file_path);
        
        let response = self.client
            .get(&url)
            .send()
            .context("Failed to send info request to backend")?;
            
        if !response.status().is_success() {
            return Err(anyhow::anyhow!(
                "Backend returned error: {} - {}",
                response.status(),
                response.text().unwrap_or_default()
            ));
        }
        
        let result = response
            .json::<MediaInfo>()
            .context("Failed to parse media info response")?;
            
        Ok(result)
    }
    
    pub fn check_health(&self) -> Result<()> {
        let url = format!("{}/health", self.base_url);
        
        let response = self.client
            .get(&url)
            .send()
            .context("Failed to send health check request")?;
            
        if response.status().is_success() {
            Ok(())
        } else {
            Err(anyhow::anyhow!("Backend health check failed"))
        }
    }
}
