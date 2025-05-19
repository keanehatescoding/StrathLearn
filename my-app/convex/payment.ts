"use node";
import { v } from "convex/values";
import { action } from "./_generated/server";
import { Pool } from "pg";

const pool = new Pool({
  connectionString: "postgresql://neondb_owner:npg_H20AghUkNFTY@ep-odd-mud-aau4cvww-pooler.westus3.azure.neon.tech/neondb?sslmode=require",
  ssl: {
    rejectUnauthorized: false
  }
});


interface PolarCustomer {
  id: string;
  email: string;
}

interface PolarSubscription {
  id: string;
  status?: string;
  recurring_interval?: string;
  current_period_end?: number | string;
  cancel_at_period_end?: boolean;
  created_at?: string;
  modified_at?: string | null;
  started_at?: string;
  ended_at?: string | null;
  ends_at?: string | null;
  canceled_at?: string | null;
  customer_id?: string;
  product_id?: string;
  discount_id?: string | null;
  checkout_id?: string;
  customer_cancellation_reason?: string | null;
  customer_cancellation_comment?: string | null;
  price_id?: string;
  user_id?: string;
  metadata?: Record<string, unknown> | null;
  amount?: number;
  currency?: string;
}


interface PolarOrderData {
  customer: PolarCustomer;
  subscription?: PolarSubscription; 
  product_id: string;
  amount: number;
  currency: string;
}


interface PolarWebhookPayload {
  type: string;
  data: PolarOrderData;
}

interface WebhookResponse {
  status: number;
  body: {
    message?: string;
    success?: boolean;
    error?: string;
    details?: string;
  };
}

async function processPolarWebhook(payload: PolarWebhookPayload): Promise<WebhookResponse> {
  try {
    console.log("Processing Polar webhook");
    console.log("Webhook type:", payload.type);
    
    if (payload.type !== "order.paid") {
      return { 
        status: 200, 
        body: { message: "Event type not handled" } 
      };
    }
    
    const orderData: PolarOrderData = payload.data;
    const customer: PolarCustomer = orderData.customer;
    const subscription: PolarSubscription | undefined = orderData.subscription;
    
    console.log(`Processing subscription for customer: ${customer.email}`);
    

    const client = await pool.connect();
    
    try {
   
      await client.query('BEGIN');
      

      const userResult = await client.query(
        'SELECT id FROM "user" WHERE email = $1',
        [customer.email]
      );
      
      let userId: string | null = null;
      if (userResult.rows.length > 0) {
        userId = userResult.rows[0].id;
        console.log(`Found user with ID: ${userId}`);
        

        await client.query(
          'UPDATE "user" SET customer_id = $1, updated_at = NOW() WHERE id = $2',
          [customer.id, userId]
        );
        console.log(`Updated user ${userId} with customer ID: ${customer.id}`);
      } else {
        console.log(`User not found with email: ${customer.email}`);
      }
      
      if (userId && subscription) {
 
        const existingSubResult = await client.query(
          'SELECT id FROM subscription WHERE id = $1',
          [subscription.id]
        );
        
        if (existingSubResult.rows.length === 0) {

          await client.query(
            `INSERT INTO subscription (
              id, user_id, status, product_id, amount, currency, interval, 
              current_period_end, cancel_at_period_end, created_at, updated_at
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())`,
            [
              subscription.id,
              userId,
              subscription.status || 'active',
              orderData.product_id,
              orderData.amount,
              orderData.currency,
              subscription.recurring_interval || 'month',
              new Date(subscription.current_period_end || (Date.now() + 30 * 24 * 60 * 60 * 1000)),
              subscription.cancel_at_period_end || false
            ]
          );
          console.log(`Added subscription ${subscription.id} to database`);
        } else {
   
          await client.query(
            `UPDATE subscription SET 
              status = $1, 
              current_period_end = $2,
              cancel_at_period_end = $3,
              updated_at = NOW()
            WHERE id = $4`,
            [
              subscription.status || 'active',
              new Date(subscription.current_period_end || (Date.now() + 30 * 24 * 60 * 60 * 1000)),
              subscription.cancel_at_period_end || false,
              subscription.id
            ]
          );
          console.log(`Updated existing subscription ${subscription.id}`);
        }
      }
      

      await client.query('COMMIT');
      console.log("Transaction committed successfully");
      
      return { 
        status: 200, 
        body: { success: true }
      };
      
    } catch (dbError) {
  
      await client.query('ROLLBACK');
      console.error("Database error:", dbError);
      
      return { 
        status: 500, 
        body: { 
          error: "Database error", 
          details: dbError instanceof Error ? dbError.message : 'Unknown error' 
        }
      };
      
    } finally {

      client.release();
    }
    
  } catch (error) {
    console.error("Webhook processing error:", error);
    
    return { 
      status: 500, 
      body: { 
        error: "Webhook processing error", 
        details: error instanceof Error ? error.message : 'Unknown error' 
      }
    };
  }
}

export const processPolarWebhookAction = action({
  args: {
    type: v.string(),
    data: v.any(),
  },
  handler: async (ctx, { type, data }) => {
   
    const orderData: PolarOrderData = {
      customer: {
        id: data.customer.id,
        email: data.customer.email,
      },
     
      ...(data.subscription ? {
        subscription: {
          id: data.subscription.id,
          status: data.subscription.status,
          recurring_interval: data.subscription.recurring_interval,
          current_period_end: data.subscription.current_period_end,
          cancel_at_period_end: data.subscription.cancel_at_period_end,
        }
      } : {}),
      product_id: data.product_id,
      amount: data.amount,
      currency: data.currency,
    };
    
    const payload: PolarWebhookPayload = { type, data: orderData };
    
    try {
      const result = await processPolarWebhook(payload);
      return result;
    } catch (error) {
      console.error("Action error:", error);
      return {
        status: 500,
        body: {
          error: "Action processing error",
          details: error instanceof Error ? error.message : 'Unknown error'
        }
      };
    }
  },
});