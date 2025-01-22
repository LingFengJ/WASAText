<script>
export default {
    data() {
        return {
            newUsername: '',
            selectedPhoto: null,
            error: null,
            success: null,
            loading: false
        }
    },
    methods: {
        async updateUsername() {
            if (!this.newUsername) return;

            this.loading = true;
            try {
                const response = await this.$axios.put('/users/me/name', 
                    {
                        name: this.newUsername
                    },
                    {
                    headers: {
                        'Authorization': `Bearer ${sessionStorage.getItem('authToken')}`,
                        'Content-Type': 'application/json'
                    },

                });

                sessionStorage.setItem('username', this.newUsername);
                this.success = 'Username updated successfully';
                this.error = null;

            } catch (error) {
                if (error.response){    
                    const statusCode = error.response.status;
                    switch (statusCode) {
                        case 400:
                            console.error('Bad request');
                            this.error = 'Bad request';
                        case 409:
                            console.error('Conflict: username already taken', error.response.data);
                            this.error = 'Username already taken';
                            break;
                        case 500:
                            console.error('Internal Server Error:', error.response.data);
                            this.error = 'Internal Server Error';
                            break;
                        default:
                            console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
                            this.error = 'Network error, failed to update username';
                    }
                } else {
                    console.error('Network error:', error);
                    this.error = 'Network error, failed to update username';
                } 
                this.success = null;
            } finally {
                this.loading = false;
            }
        },

        async updatePhoto() {
            if (!this.selectedPhoto) return;

            const formData = new FormData();
            formData.append('photo', this.selectedPhoto);

            try {
                const response = await this.$axios.put('/users/me/photo', 
                formData,
                {
                    headers: {
                        'Authorization': `Bearer ${sessionStorage.getItem('authToken')}`
                    }
                });

                this.success = 'Profile photo updated successfully';
                this.error = null;

            } catch (error) {
                console.error('Error updating photo:', error);
                this.error = 'Failed to update profile photo';
                this.success = null;
            }
        },

        onPhotoSelected(event) {
            this.selectedPhoto = event.target.files[0];
            if (this.selectedPhoto) {
                this.updatePhoto();
            }
        }
    }
}
</script>

<template>
    <div class="user-profile p-4">
        <h2>Profile Settings</h2>

        <div v-if="error" class="alert alert-danger">
            {{ error }}
        </div>

        <div v-if="success" class="alert alert-success">
            {{ success }}
        </div>

        <div class="mb-4">
            <h4>Change Username</h4>
            <form @submit.prevent="updateUsername" class="mb-3">
                <div class="input-group">
                    <input type="text" 
                           class="form-control" 
                           v-model="newUsername" 
                           placeholder="New username"
                           required>
                    <button type="submit" 
                            class="btn btn-primary"
                            :disabled="loading">
                        {{ loading ? 'Updating...' : 'Update' }}
                    </button>
                </div>
            </form>
        </div>

        <div class="mb-4">
            <h4>Profile Photo</h4>
            <input type="file" 
                   class="form-control" 
                   accept="image/*"
                   @change="onPhotoSelected">
        </div>
    </div>
</template>