// Supabase Edge Function: Send Push Notification
import { serve } from 'https://deno.land/std@0.168.0/http/server.ts';
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2';

serve(async (req) => {
  try {
    const { user_id, title, body, data } = await req.json();

    const supabase = createClient(
      Deno.env.get('SUPABASE_URL')!,
      Deno.env.get('SUPABASE_SERVICE_ROLE_KEY')!
    );

    // Get user's FCM token
    const { data: user } = await supabase
      .from('users')
      .select('fcm_token')
      .eq('user_id', user_id)
      .single();

    if (!user?.fcm_token) {
      return new Response(JSON.stringify({ error: 'No FCM token found' }), {
        status: 404,
        headers: { 'Content-Type': 'application/json' },
      });
    }

    // Send notification via FCM
    const fcmResponse = await fetch('https://fcm.googleapis.com/fcm/send', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `key=${Deno.env.get('FCM_SERVER_KEY')}`,
      },
      body: JSON.stringify({
        to: user.fcm_token,
        notification: { title, body },
        data: data || {},
      }),
    });

    const fcmResult = await fcmResponse.json();

    return new Response(JSON.stringify({ success: true, fcm_result: fcmResult }), {
      headers: { 'Content-Type': 'application/json' },
    });
  } catch (error) {
    return new Response(JSON.stringify({ error: error.message }), {
      status: 500,
      headers: { 'Content-Type': 'application/json' },
    });
  }
});
