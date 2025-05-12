import { betterAuth } from "better-auth";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import db from "../db/index.js";
import * as schema from "../db/schema.js";
import { jwt } from "better-auth/plugins"

export const auth = betterAuth({
    database: drizzleAdapter(db, {
        provider: "pg",
        schema: schema  
    }),
    plugins:[
        jwt(), 

    ],
    
    emailAndPassword: {
        enabled: true,
        normalizeEmail: true
    },
    socialProviders: { 
        github: { 
            clientId: "Ov23lipI60Uz3L3D36Us", 
            clientSecret: "f7020ec87447f66e1da60b1e1758e885bdda7bff", 
        } 
    }, 
    trustedOrigins: [
        'http://localhost:3000',
        'http://localhost:5173',
        'https://api.singularity.co.ke',
    ]
});