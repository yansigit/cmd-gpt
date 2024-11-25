use clap::Parser;
use reqwest;
use serde::{Deserialize, Serialize};
use std::error::Error;
use dotenv::dotenv;

#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Args {
    #[arg(default_value = "")]
    prompt: String,

    #[arg(long, short)]
    generate_powershell_command: Option<String>,

    #[arg(long, short, default_value = "claude-3-haiku")]
    model: String,
}

#[derive(Serialize, Debug)]
struct OpenAIChatRequest {
    model: String,
    messages: Vec<Message>,
    temperature: f32,
}

#[derive(Serialize, Debug)]
struct AnthropicChatRequest {
    model: String,
    messages: Vec<Message>,
    max_tokens: u32,
}

#[derive(Serialize, Debug)]
struct Message {
    role: String,
    content: String,
}

#[derive(Deserialize, Debug)]
struct OpenAIChatResponse {
    id: String,
    object: String,
    created: u64,
    model: String,
    choices: Vec<OpenAIChoice>,
    usage: Usage,
}

#[derive(Deserialize, Debug)]
struct AnthropicChatResponse {
    id: String,
    content: Vec<AnthropicContent>,
}

#[derive(Deserialize, Debug)]
struct AnthropicContent {
    text: String,
}

#[derive(Deserialize, Debug)]
struct OpenAIChoice {
    index: u32,
    message: ResponseMessage,
    finish_reason: String,
}

#[derive(Deserialize, Debug)]
struct ResponseMessage {
    role: String,
    content: String,
}

#[derive(Deserialize, Debug)]
struct Usage {
    prompt_tokens: u32,
    completion_tokens: u32,
    total_tokens: u32,
}

async fn call_openai(prompt: &str, model: &str, api_key: &str) -> Result<String, Box<dyn Error>> {
    let client = reqwest::Client::new();
    let request_body = OpenAIChatRequest {
        model: model.to_string(),
        messages: vec![Message {
            role: "user".to_string(),
            content: prompt.to_string(),
        }],
        temperature: 0.7,
    };

    let response = client
        .post("https://api.openai.com/v1/chat/completions")
        .header("Authorization", format!("Bearer {}", api_key))
        .header("Content-Type", "application/json")
        .json(&request_body)
        .send()
        .await?;

    if !response.status().is_success() {
        let error_text = response.text().await?;
        return Err(error_text.into());
    }

    let response_data = response.json::<OpenAIChatResponse>().await?;
    Ok(response_data.choices[0].message.content.clone())
}

async fn call_anthropic(prompt: &str, model: &str, api_key: &str) -> Result<String, Box<dyn Error>> {
    let client = reqwest::Client::new();
    let request_body = AnthropicChatRequest {
        model: model.to_string(),
        messages: vec![Message {
            role: "user".to_string(),
            content: prompt.to_string(),
        }],
        max_tokens: 1024,
    };

    let response = client
        .post("https://api.anthropic.com/v1/messages")
        .header("x-api-key", api_key)
        .header("anthropic-version", "2023-06-01")
        .header("Content-Type", "application/json")
        .json(&request_body)
        .send()
        .await?;

    if !response.status().is_success() {
        let error_text = response.text().await?;
        return Err(error_text.into());
    }

    let response_data = response.json::<AnthropicChatResponse>().await?;
    Ok(response_data.content[0].text.clone())
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let args = Args::parse();
    dotenv().ok();
    
    let prompt = if let Some(cmd_prompt) = args.generate_powershell_command {
        format!("Generate a PowerShell command for the following task: {}", cmd_prompt)
    } else {
        args.prompt
    };

    let response = match args.model.as_str() {
        "gpt-3.5-turbo" | "gpt-4" => {
            let api_key = std::env::var("OPENAI_API_KEY")
                .expect("OPENAI_API_KEY environment variable is required");
            call_openai(&prompt, &args.model, &api_key).await?
        }
        "claude-3-haiku" | "claude-3-opus" | "claude-3-sonnet" => {
            let api_key = std::env::var("ANTHROPIC_API_KEY")
                .expect("ANTHROPIC_API_KEY environment variable is required");
            call_anthropic(&prompt, &args.model, &api_key).await?
        }
        _ => return Err("Unsupported model".into()),
    };

    println!("{}", response);
    Ok(())
}
