<script>
export default {
    data() {
        return {
            username: '',
            password: '',
            error: null
        }
    },
    methods: {
        async login() {
            try {
                console.log('Sending login request with:', {
                    username: this.username,
                    password: this.password
                });

                const response = await fetch('http://localhost:3000/session', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        name: this.username,
                        password: this.password
                    })
                });
                
                console.log('Response status:', response.status);
                const responseText = await response.text();
                console.log('Response text:', responseText);

                if (response.ok) {
                    const data = JSON.parse(responseText);
                    console.log('Login successful, got identifier:', data.identifier);
                    sessionStorage.setItem('authToken', data.identifier);
                    sessionStorage.setItem('username', this.username);
                    this.$router.push('/conversations');
                } else {
                    try {
                        const errorData = JSON.parse(responseText);
                        this.error = errorData.message || 'Login failed';
                    } catch (e) {
                        this.error = `Login failed - ${responseText}`;
                    }
                }
            } catch (error) {
                console.error('Login error:', error);
                console.error('Error details:', {
                    message: error.message,
                    stack: error.stack
                });
                this.error = `Login failed - ${error.message}`;
            }
        }
    }
}
</script>

<template>
    <div class="login-form p-4">
        <h2 class="mb-4">Login to WASAText</h2>
        <div v-if="error" class="alert alert-danger">{{ error }}</div>
        
        <form @submit.prevent="login" class="needs-validation">
            <div class="mb-3">
                <label for="username" class="form-label">Username:</label>
                <input 
                    type="text" 
                    id="username"
                    v-model="username" 
                    class="form-control" 
                    required
                    minlength="3"
                    maxlength="16"
                >
            </div>
            
            <div class="mb-3">
                <label for="password" class="form-label">Password:</label>
                <input 
                    type="password" 
                    id="password"
                    v-model="password" 
                    class="form-control" 
                    required
                    minlength="4"
                >
            </div>
            
            <button type="submit" class="btn btn-primary">Login</button>
        </form>
    </div>
</template>

<style scoped>
.login-form {
    max-width: 400px;
    margin: 0 auto;
}
</style>