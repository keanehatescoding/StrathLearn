import { httpAction } from "./_generated/server";
import { httpRouter } from "convex/server";
import { api } from "./_generated/api";

const http = httpRouter();

export const polarWebhook = httpAction(async (ctx, request) => {
  try {
    const payload = await request.json();
    
    const result = await ctx.runAction(api.payment.processPolarWebhookAction, payload);
    
    return new Response(JSON.stringify(result.body), {
      status: result.status,
      headers: { "Content-Type": "application/json" }
    });
  } catch (error) {
    return new Response(JSON.stringify({
      error: "Error processing webhook",
      details: error instanceof Error ? error.message : 'Unknown error'
    }), {
      status: 500,
      headers: { "Content-Type": "application/json" }
    });
  }
});

http.route({
  path: "/webhook/polar",
  method: "POST",
  handler: polarWebhook,
});

export default http;