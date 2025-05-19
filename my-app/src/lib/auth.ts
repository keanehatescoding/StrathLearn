import { betterAuth } from "better-auth";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { jwt } from "better-auth/plugins";
import { polar } from "@polar-sh/better-auth"; 
import { Polar } from "@polar-sh/sdk"; 
import db from "../db/index.js";
import * as schema from "../db/schema.js";
import { eq } from "drizzle-orm";

const client = new Polar({
    accessToken: "polar_oat_2OkiVqy3ugQfOJtq8d077u7WM73ZCgmNAVrwR3Xnruh", 
    server: 'sandbox',
});


export const auth = betterAuth({
    database: drizzleAdapter(db, {
        provider: "pg",
        schema: schema  
    }),
    plugins: [
        jwt(), 
        polar({
            client,
            createCustomerOnSignUp: true,
            enableCustomerPortal: true,
            checkout: {
                enabled: true,
                products: [
                    {
                        productId: "9398d8de-c3c3-423e-b515-52f8b5f67596", 
                        slug: "course"
                    }
                ],
                successUrl: "/success?checkout_id={CHECKOUT_ID}",
                authenticatedUsersOnly: true,
            },
            webhooks: {
                secret: "3f6f2f49d07f4d02997af03135389a5f",
              onOrderPaid: async (order) => {
    try {
        console.log("Received order paid webhook!");
        
        // Try multiple ways to find the user
        const userEmail = order.data.customer.email;
        const externalId = order.data.customer.externalId;
        
        let user = null;
        
        // First try by email
        if (userEmail) {
            user = await db.query.user.findFirst({
                where: eq(schema.user.email, userEmail)
            });
        }
        
        // If not found and we have externalId, try that
        if (!user && externalId) {
            user = await db.query.user.findFirst({
                where: eq(schema.user.id, externalId)
            });
        }
        
        if (!user) {
            console.error(`Could not find user with email: ${userEmail} or id: ${externalId}`);
            // You might want to store this information somewhere for manual resolution
        } else {
            console.log(`Found user: ${user.id}`);
            
            // Get subscription data
            const subscriptionData = order.data.subscription;
            
            // Insert subscription with the actual user ID
            await db.insert(schema.subscription).values({
                id: subscriptionData?.id || `sub_manual_${Date.now()}`,
                userId: user.id,
                status: order.data.status || "active",
                productId: order.data.productId || "9398d8de-c3c3-423e-b515-52f8b5f67596",
                amount: order.data.amount || 1000,
                currency: order.data.currency || "usd",
                interval: subscriptionData?.recurringInterval || "month",
                currentPeriodEnd: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000),
                cancelAtPeriodEnd: false,
                createdAt: new Date(),
                updatedAt: new Date()
            });
            
            console.log("Added subscription record to database!");
            
            // Update user's customerId if needed
            if (!user.customerId || user.customerId !== order.data.customer.id) {
                await db.update(schema.user)
                    .set({ customerId: order.data.customer.id })
                    .where(eq(schema.user.id, user.id));
                console.log(`Updated user ${user.id} with customer ID ${order.data.customer.id}`);
            }
        }
    } catch (error) {
        console.error("Error in webhook handler:", error);
        // Consider implementing retry logic or notification system
    }
},
                onPayload: async (payload) => {
                    console.log("Received webhook payload type:", payload.type);
                }
            }
        })
    ],
    emailAndPassword: {
        enabled: true,
    },
    socialProviders: { 
        github: { 
            clientId: process.env.GITHUB_CLIENT_ID as string, 
            clientSecret: process.env.GITHUB_CLIENT_SECRET as string, 
        } 
    }, 
    trustedOrigins: [
        'http://localhost:3000',
        'http://localhost:5173',
        'https://api.singularity.co.ke',
        "https://ed1e-196-216-95-131.ngrok-free.app",
        "https://5477-196-216-95-131.ngrok-free.app",
        "https://ea63-102-0-21-150.ngrok-free.app"
    ]
});