// Raw localStorage to avoid module evaluation issues
if (typeof localStorage !== 'undefined') {
  localStorage.setItem('accessToken', 'dummy-token')
}
