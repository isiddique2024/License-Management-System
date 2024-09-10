import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import crypto from "crypto";

function generateNonce() {
  return crypto.randomBytes(16).toString("base64");
}

const nonce = generateNonce();

export default defineConfig({
  base: "/",
  plugins: [
    react(),
    {
      name: "html-inject-data-preload-attr",
      enforce: "post",
      transformIndexHtml(html) {
        const regex = /<(link|style|script)/gi;
        const replacement = '<$1 data-preload="true"';
        return html.replace(regex, replacement);
      },
    },
    {
      name: "html-nonce-injector",
      enforce: "post",
      transformIndexHtml(html) {
        return html.replace(/<script/gi, `<script nonce="${nonce}"`);
      },
    },
  ],
  server: {
    open: false,
    port: 3000,
    strictPort: true,
    host: true,
    fs: {
      strict: true,
    },
    watch: {
      usePolling: true,
    },
    headers: {
      "Content-Security-Policy": `default-src 'self'; script-src 'self' 'nonce-${nonce}'; style-src 'self' 'sha256-Erd4Drq6hrHg1gjcimFFVdyc+XvYff4xrcySRGpBUT8=' 'sha256-cMxjrcQcx7TgbXGafMFTYrg5FHNIefSegDbcQrG/XbQ=' 'sha256-2J6GAaw6QU8MIh5BMascsHgZS322IWlb4jXz49SL+f4='; img-src 'self' data: http://www.w3.org/; connect-src 'self'; font-src 'self'; object-src 'none'; base-uri 'self'; form-action 'self';`,
    },
  },
  build: {
    outDir: "dist",
    sourcemap: false,
    target: "esnext",
    minify: "terser", // Use Terser for minification
    cssCodeSplit: true,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ["react", "react-dom"],
        },
      },
    },
    terserOptions: {
      compress: {
        drop_console: true,
      },
    },
  },
  optimizeDeps: {
    include: ["react", "react-dom"],
    exclude: ["@vitejs/plugin-react"],
  },
});
