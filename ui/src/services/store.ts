import Backend from "@/services/backend";
// count the number of hits from all backends
let hits = 0;

// keep backends here
const backends = new Array<Backend>();

// Global store that can be accessed from anywhere
export default {
  get backends() {
    return backends;
  },
  set backends(value) {
    // clear backends
    backends.splice(0, backends.length);
    // add new backends
    backends.push(...value);
  },
  incrementHits: (c = 1) => {
    hits += c;
  },
  get hits() {
    return hits;
  },
  set hits(value) {
    hits = value;
  },
};
