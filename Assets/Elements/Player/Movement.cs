using System.Collections;
using System.Collections.Generic;
using System;
using UnityEngine;

public class Movement : MonoBehaviour {
    public float timeout = 1f / 24f;
    public float speed;
    public float bulletForce;
    public GameObject bullet;

    private float updateTimeout;
    private AudioSource audio;
    private RemoteManager remote;
    private Rigidbody2D rb;
    private float dir;

    void Start() {
        RenderSettings.ambientLight = Color.white;
        RenderSettings.ambientIntensity = 100f;

        updateTimeout = timeout;
        rb = GetComponent<Rigidbody2D>();
        remote = GetComponent<RemoteManager>();
        audio = GetComponent<AudioSource>();
    }

    void Update() {
        var dir = GetDirection();

        // Movement
        if (Input.GetMouseButton(0)) {
            this.dir = Mathf.Rad2Deg * Mathf.Atan2(dir.y, dir.x);
            gameObject.transform.rotation = GetRotation(this.dir);
            if (!audio.isPlaying) audio.Play();
            rb.velocity = dir * speed;
        } else {
            rb.velocity = Vector2.zero;
        }

        // Shooting
        if (Input.GetKeyDown(KeyCode.Space)) {
            var bullet = Instantiate(this.bullet, gameObject.transform.position, Quaternion.identity);
            var rbBullet = bullet.GetComponent<Rigidbody2D>();
            rbBullet.AddForce(dir * bulletForce);
        }
    }

    void FixedUpdate() {
        updateTimeout = Mathf.Max(0f, updateTimeout - Time.deltaTime);
        if (updateTimeout == 0 && remote != null && gameObject.transform.hasChanged) {
            remote.UpdateSelf(dir, new Vector2(gameObject.transform.position.x, gameObject.transform.position.y));
            updateTimeout = timeout;
        }
    }

    /* --- Utils --- */
    Vector2 GetDirection() {
        var halfScreen = new Vector2(Screen.width, Screen.height) / 2;
        var position = new Vector2(Input.mousePosition.x, Input.mousePosition.y) - halfScreen;
        return position.normalized;
    }

    Quaternion GetRotation (float deg) {
        var rotation = Quaternion.AngleAxis(deg, Vector3.forward);
        return rotation * Quaternion.AngleAxis(-deg, Vector3.right);
    }
}
