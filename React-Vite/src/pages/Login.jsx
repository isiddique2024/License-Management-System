import React from "react";
import { signal } from "@preact/signals-react";
const loading = signal(false);
const username = signal("");
const password = signal("");
const error = signal(null);

const LoginPage = () => {
  const handleLogin = async (event) => {
    event.preventDefault();
    loading.value = true;
    error.value = null;
  };

  return (
    <div className="max-w-lg flex-auto flex-col items-center justify-center mx-auto mt-16">
      <div className="main-content flex-auto">
        <div className="container mx-auto p-4 bg-[#333b45] flex-auto flex-col justify-center items-center rounded-lg shadow-lg w-full">
          <h1 className="heading text-white text mb-4">Login</h1>
          <form onSubmit={handleLogin} className="w-full">
            <div className="mb-4 flex flex-col items-center">
              <label className="label text-white" htmlFor="username">
                Username/Email
              </label>
              <input
                type="text"
                id="username"
                className="input p-2 rounded"
                value={username.value}
                onChange={(e) => (username.value = e.target.value)}
                required
              />
            </div>
            <div className="mb-4 flex flex-col items-center">
              <label className="label text-white" htmlFor="password">
                Password
              </label>
              <input
                type="password"
                id="password"
                className="input p-2 rounded"
                value={password.value}
                onChange={(e) => (password.value = e.target.value)}
                required
              />
            </div>
            {error.value && (
              <div className="mb-4 text-red-500 text-sm">{error.value}</div>
            )}
            <div className="mb-4 flex items-center justify-center space-x-3">
              <button
                type="submit"
                className="button flex-shrink p-2 rounded bg-blue-500 text-white"
              >
                Login
              </button>
              {loading.value && (
                <span className="loading loading-spinner loading-lg"></span>
              )}
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
