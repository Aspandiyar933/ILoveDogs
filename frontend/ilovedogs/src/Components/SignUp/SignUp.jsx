import React, { useState } from 'react';

const SignupForm = () => {
  const [loading, setLoading] = useState(false); // Added loading state

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true); // Set loading state to true when form is submitted

    try {
      const formData = new FormData(e.currentTarget); // Create FormData object from form
      const response = await fetch('http://localhost:3000/api/v1/users/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'multipart/form-data', // Ensure correct Content-Type
        },
        body: formData, // Send FormData instead of JSON
      });

      if (!response.ok) {
        throw new Error('Failed to sign up');
      }

      // Handle successful signup
      console.log('User signed up successfully');
    } catch (error) {
      console.error('Error signing up:', error.message);
    } finally {
      setLoading(false); // Reset loading state after request is complete
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label htmlFor="name">Name:</label>
        <input type="text" id="name" name="name" />
      </div>
      <div>
        <label htmlFor="username">Username:</label>
        <input type="text" id="username" name="username" />
      </div>
      <div>
        <label htmlFor="email">Email:</label>
        <input type="email" id="email" name="email" />
      </div>
      <div>
        <label htmlFor="password">Password:</label>
        <input type="password" id="password" name="password" />
      </div>
      <button type="submit">Sign Up{loading && '...'}</button> {/* Show loading indicator */}
    </form>
  );
};

export default SignupForm;
