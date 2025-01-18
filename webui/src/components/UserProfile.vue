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
                const response = await fetch('http://localhost:3000/users/me/name', {
                    method: 'PUT',
                    headers: {
                        'Authorization': `Bearer ${sessionStorage.getItem('authToken')}`,
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        name: this.newUsername
                    })
                });

                if (response.ok) {
                    sessionStorage.setItem('username', this.newUsername);
                    this.success = 'Username updated successfully';
                    this.error = null;
                } else {
                    this.error = 'Failed to update username';
                    this.success = null;
                }
            } catch (error) {
                this.error = 'Network error';
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
                const response = await fetch('http://localhost:3000/users/me/photo', {
                    method: 'PUT',
                    headers: {
                        'Authorization': `Bearer ${sessionStorage.getItem('authToken')}`
                    },
                    body: formData
                });

                if (response.ok) {
                    this.success = 'Profile photo updated successfully';
                    this.error = null;
                } else {
                    this.error = 'Failed to update profile photo';
                    this.success = null;
                }
            } catch (error) {
                this.error = 'Network error';
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