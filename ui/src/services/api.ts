import Backend from "./backend";
import Certificate from "./certificate";
import Hit from "./hit";

/**
 * Api Class to handle the api calls.
 */
export default class Api {
  BASE_URL = "http://localhost:8080/api/v1";

  // Get the list of backend
  async getBackends(): Promise<Backend[]> {
    return fetch(`${this.BASE_URL}/servers`).then((res) => res.json());
  }

  // Get certificate information for a backend
  async getCert(servername: string): Promise<Certificate> {
    return fetch(`${this.BASE_URL}/cert/${servername}`).then(
      (res) => res.json(),
      () => {
        return Promise.reject("No certificate found");
      }
    );
  }

  // Get the state (up/down) of a backend
  async getBackendState(servername: string): Promise<boolean> {
    return fetch(`${this.BASE_URL}/state/${servername}`).then((res) =>
      res.json()
    );
  }

  // Get the memory allocation size of the proxy
  async getServerMemAlloc(): Promise<number> {
    return fetch(`${this.BASE_URL}/runtime/mem/alloc`).then((res) =>
      res.json()
    );
  }

  // Get the proxy version
  async getServerVersion(): Promise<string> {
    return fetch(`${this.BASE_URL}/version`).then((res) => res.json());
  }

  // Get stats about one backend
  async getBackendStats(servername: string): Promise<Array<Hit>> {
    return fetch(`${this.BASE_URL}/stats/${servername}`).then((res) =>
      res.json()
    );
  }

  // Change state of the backend
  async setBackend(name: string, _new: Backend) {
    return fetch(`${this.BASE_URL}/backend/${name}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(_new),
    });
  }
}
