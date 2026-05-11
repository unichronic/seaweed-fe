import { initializeApp, getApps, FirebaseApp } from 'firebase/app';
import { getAuth, Auth, User } from 'firebase/auth';

const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY || "fallback-key",
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
  storageBucket: process.env.NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET,
  messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
  appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID,
};

let appInstance: FirebaseApp | undefined;
let authInstance: Auth;

if (process.env.NEXT_PUBLIC_AUTH_BYPASS === 'true') {
  const localUser = {
    uid: 'local-user-123',
    email: 'local@example.com',
    getIdToken: async () => 'dummy-token',
  } as unknown as User;

  authInstance = {
    currentUser: localUser,
  } as unknown as Auth;
} else {
  appInstance = getApps().length === 0 ? initializeApp(firebaseConfig) : getApps()[0];
  authInstance = getAuth(appInstance);
}

export const auth = authInstance;
export default appInstance;
