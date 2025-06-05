import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      // Proxying all requests starting with /api to localhost:8080
      "/api": {
        target: "http://localhost:8080",
        // changeOrigin: true, // Needed for virtual hosted sites
        // rewrite: (path) => path.replace(/^\/api/, ''), // Optional: remove /api prefix from the URL sent to the backend
      },
      // You can add more proxy rules if needed, for example:
      // '/auth': {
      //   target: 'http://localhost:8080',
      //   changeOrigin: true,
      // },
    },
  },
});
