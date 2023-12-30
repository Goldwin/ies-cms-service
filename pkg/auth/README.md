# Authentication and Authorization Service



Registration Flow:
1. User enter email
2. Send OTP to email
3. User enter OTP
4. Account Created, but not linked. Send auth token
5. User Enter First Name And Last Name
6. Fetch then select existing Profile, or Create New Person Profile
7. Link account with Person Profile

Authentication Flow (username/password):
1. User enter username & password
2. receive auth token

Authentication Flow (email/otp):
1. User enter email
2. Send OTP to email
3. User enter OTP
4. Send Auth token