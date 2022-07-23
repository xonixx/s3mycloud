function backendBaseUrl(): string {
  if (location.href.indexOf("http://localhost:3000") === 0) {
    return "http://127.0.0.1:8080";
  }
  throw Error("Can't determine backend base URL");
}

export const BACKEND_BASE_URL: string = backendBaseUrl();
