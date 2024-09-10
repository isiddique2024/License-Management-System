// import Keycloak, { KeycloakOnLoad } from "keycloak-js";
// import { signal } from "@preact/signals-react";

// const keycloakConfig = {
//   url: `https://localhost/`,
//   realm: "demo",
//   clientId: "real-client",
// };

// const keycloak = new Keycloak(keycloakConfig);
// const keycloakInitialized = signal<boolean>(false);
// let initPromise: Promise<void> | null = null;

// export const initKeycloak = (
//   onLoad: KeycloakOnLoad = "login-required",
//   redirectUri = `https://localhost/dashboard` // ${window.location.origin}/
// ) => {
//   if (!initPromise) {
//     initPromise = keycloak
//       .init({
//         onLoad,
//         checkLoginIframe: false,
//         checkLoginIframeInterval: 25,
//         redirectUri: redirectUri,
//       })
//       .then((authenticated) => {
//         keycloakInitialized.value = true;
//         if (!authenticated) {
//           console.warn("User is not authenticated");
//         }
//       })
//       .catch((err) => {
//         keycloakInitialized.value = false;
//         console.error(err);
//       });
//   }
//   return initPromise;
// };

// export default keycloak;
// export { keycloakInitialized };
import Keycloak, { KeycloakOnLoad } from "keycloak-js";
import { signal } from "@preact/signals-react";

// Access environment variables using Vite's import.meta.env
const keycloakConfig = {
  url: import.meta.env.VITE_PUBLIC_URL, // Default to localhost if env variable is not set
  realm: import.meta.env.VITE_REALM, // Use env variable or default to "demo"
  clientId: import.meta.env.VITE_USER_CLIENT_ID, // Use env variable or default to "real-client"
};

const keycloak = new Keycloak(keycloakConfig);
const keycloakInitialized = signal<boolean>(false);
let initPromise: Promise<void> | null = null;

export const initKeycloak = (
  onLoad: KeycloakOnLoad = "login-required",
  redirectUri = `${import.meta.env.VITE_PUBLIC_URL}dashboard` // Use env variable or default
) => {
  if (!initPromise) {
    initPromise = keycloak
      .init({
        onLoad,
        checkLoginIframe: false,
        checkLoginIframeInterval: 25,
        redirectUri: redirectUri,
      })
      .then((authenticated) => {
        keycloakInitialized.value = true;
        if (!authenticated) {
          console.warn("User is not authenticated");
        }
      })
      .catch((err) => {
        keycloakInitialized.value = false;
        console.error(err);
      });
  }
  return initPromise;
};

export default keycloak;
export { keycloakInitialized };
