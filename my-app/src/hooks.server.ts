import { auth } from "$lib/auth.js";
import { svelteKitHandler } from "better-auth/svelte-kit";
import { Polar } from "@polar-sh/sdk";
import type { Handle } from "@sveltejs/kit";
import { redirect } from "@sveltejs/kit";

const polar = new Polar({
  accessToken: process.env.POLAR_ACCESS_TOKEN || "polar_oat_2OkiVqy3ugQfOJtq8d077u7WM73ZCgmNAVrwR3Xnruh",
  server: 'sandbox',
});

export const handle: Handle = async ({ event, resolve }) => {
  const path = event.url.pathname;
  console.log("Current path:", path);
  console.log("Full URL:", event.url.toString());
  

  if (path === '/success' && event.url.searchParams.has('customer_session_token')) {
    const token = event.url.searchParams.get('customer_session_token');
    console.log("Got Polar token:", token);
   
    return await resolve(event);
  }
  
  const publicRoutes = ['/signin', '/signup', '/forgot-password', '/api', '/success'];
  
  if (publicRoutes.some(route => path.startsWith(route))) {
    return await svelteKitHandler({ event, resolve, auth });
  }
  
  let session;
  try {
    session = await auth.api.getSession({
      headers: event.request.headers
    });
  } catch (err) {
    console.error("Error getting session:", err);
    if (path !== '/') {
      return redirect(303, '/');
    }
  }
  
  const isAuthenticated = !!session?.user;
  
  if (!isAuthenticated && path !== '/') {
    return redirect(303, '/');
  }
  
  if (isAuthenticated) {
    try {
        const userId = session?.user.id;
        if (!userId) {
            throw new Error("User ID not found in session");
        } 
        const customerState = await polar.customers.getStateExternal({
            externalId: userId,
        });
        
        const hasPaidSubscription = customerState.activeSubscriptions && 
                                   customerState.activeSubscriptions.length > 0 && 
                                   customerState.activeSubscriptions[0].status === 'active';
        
        if (!hasPaidSubscription) {
            return redirect(303, '/api/auth/checkout/course');
        }
        
        // Now check paths after verifying subscription
        if (path === '/api/auth/checkout/course' || path === '/challenge') {
            return await svelteKitHandler({ event, resolve, auth });
        }
    } catch (polarError) {
        console.error("Error getting Polar subscription status:", polarError);
        return redirect(303, '/api/auth/checkout/course');
    }
  }
  

  return await svelteKitHandler({ event, resolve, auth });
};