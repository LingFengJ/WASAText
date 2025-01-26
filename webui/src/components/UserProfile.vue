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
                    }
                );

                sessionStorage.setItem('username', this.newUsername);
                this.success = 'Username updated successfully';
                this.error = null;
                location.reload();
            } catch (error) {
                if (error.response){    
                    const statusCode = error.response.status;
                    switch (statusCode) {
                        case 400:
                            console.error('Bad request');
                            this.error = 'Bad request';
                            break;  // Added missing break
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

            try {
                // Read the file as binary data
                const reader = new FileReader();
                reader.onload = async (e) => {
                    try {
                        const response = await this.$axios.put('/users/me/photo', 
                            e.target.result,  // Send binary data directly
                            {
                                headers: {
                                    'Authorization': `Bearer ${sessionStorage.getItem('authToken')}`,
                                    'Content-Type': this.selectedPhoto.type  // Use the file's MIME type
                                }
                            }
                        );

                        this.success = 'Profile photo updated successfully';
                        this.error = null;
                    } catch (error) {
                        console.error('Error updating photo:', error);
                        this.error = error.response?.data || 'Failed to update profile photo';
                        this.success = null;
                    }
                };

                reader.onerror = () => {
                    this.error = 'Failed to read the selected file';
                    this.success = null;
                };

                // Read the file as binary data
                reader.readAsArrayBuffer(this.selectedPhoto);

            } catch (error) {
                console.error('Error handling photo:', error);
                this.error = 'Failed to process the selected file';
                this.success = null;
            }
        },

        onPhotoSelected(event) {
            const file = event.target.files[0];
            if (file) {
                if (!file.type.startsWith('image/')) {
                    this.error = 'Please select an image file';
                    this.success = null;
                    return;
                }
                this.selectedPhoto = file;
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