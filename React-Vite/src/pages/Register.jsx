import React from "react";
import { signal } from "@preact/signals-react";
const loading = signal(false);
const error = signal();
const username = signal("");
const password = signal("");
const email = signal("");
const confirmPassword = signal("");
const RegisterPage = () => {
  const handleSubmit = async (e) => {
    e.preventDefault();
  };

  return (
    <div className="max-w-lg flex-auto flex-col items-center justify-center">
      <div className="main-content flex-auto">
        <div className="container mx-auto p-4 bg-[#333b45] flex-auto flex-col justify-center items-center rounded-lg shadow-lg w-full">
          <h1 className="heading">Register</h1>
          <form onSubmit={handleSubmit} className="w-full">
            <div className="mb-4 flex flex-col items-center">
              <label className="label" htmlFor="username">
                Username
              </label>
              <input
                type="text"
                id="username"
                className="input"
                value={username}
                onChange={(e) => (username.value = e.target.value)}
                required
              />
            </div>
            <div className="mb-4 flex flex-col items-center">
              <label className="label" htmlFor="email">
                Email
              </label>
              <input
                type="email"
                id="email"
                className="input"
                value={email}
                onChange={(e) => (email.value = e.target.value)}
                required
              />
            </div>
            <div className="mb-4 flex flex-col items-center">
              <label className="label" htmlFor="password">
                Password
              </label>
              <input
                type="password"
                id="password"
                className="input"
                value={password}
                onChange={(e) => (password.value = e.target.value)}
                required
              />
            </div>
            <div className="mb-4 flex flex-col items-center">
              <label className="label" htmlFor="confirmPassword">
                Confirm Password
              </label>
              <input
                type="password"
                id="confirmPassword"
                className="input"
                value={confirmPassword}
                onChange={(e) => (confirmPassword.value = e.target.value)}
                required
              />
            </div>
            {error && <div className="mb-4 text-red-500 text-sm">{error}</div>}
            <div className="mb-4 flex items-center justify-center space-x-3">
              <button type="submit" className="button flex-shrink">
                Register
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

export default RegisterPage;
